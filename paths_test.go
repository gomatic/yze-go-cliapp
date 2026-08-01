package pkgstd

// White-box tests for the package-classification rules. They decide WHETHER a
// package is judged at all, so a false positive imposes command-package
// obligations on a helper that has none, and a false negative exempts a real
// command from every check in this analyzer — silently, with a green gate.

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

// TestIsCommandPackageMatchesExactlyOneSegmentBelowTheMarker names
// isCommandPackage's claim. "Exactly one path segment follows the marker" is
// the whole rule: a helper nested beneath a command
// (.../commands/greet/internal/render) is ordinary code and carries none of the
// per-package obligations, while the command package itself carries all of
// them. Matching a prefix instead would impose an entry-point requirement on
// every helper; matching too narrowly would exempt the command.
func TestIsCommandPackageMatchesExactlyOneSegmentBelowTheMarker(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path importPath
		why  string
		want bool
	}{
		{path: "m/internal/app/commands/greet", want: true, why: "exactly one segment below the marker"},
		{path: "example.com/x/internal/app/commands/serve", want: true, why: "the module prefix is irrelevant"},

		{path: "m/internal/app/commands/greet/internal/render", want: false, why: "a nested helper is not a command"},
		{path: "m/internal/app/commands/greet/sub", want: false, why: "any deeper segment disqualifies"},
		{path: "m/internal/app/commands", want: false, why: "the marker directory itself holds no command"},
		{path: "m/internal/app/commands/", want: false, why: "an empty segment names no command"},
		{path: "m/internal/app/command/greet", want: false, why: "the marker is 'commands', not 'command'"},
		{path: "m/internal/commands/greet", want: false, why: "the marker is the full app/commands path"},
		{path: "m/pkg/greet", want: false, why: "an unrelated package"},
		{path: "", want: false, why: "an empty path"},
	} {
		assert.Equal(t, tc.want, isCommandPackage(tc.path), "isCommandPackage(%q): %s", tc.path, tc.why)
	}
}

// TestIsScaffoldingPackageSkipsTheDriverSynthesizedPasses names the guard run
// depends on. The analysis driver synthesizes an external test package and a
// test-main package for any package that has tests, and NEITHER carries the
// command source — so each would trip the missing-entry-point check and report
// a violation against a package that does not exist on disk. Skipping them is
// also what guarantees pass.Files is non-empty, which is why checkCommandFunc
// may index it without a bounds check.
func TestIsScaffoldingPackageSkipsTheDriverSynthesizedPasses(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		path string
		why  string
		want bool
	}{
		{name: "greet_test", path: "m/internal/app/commands/greet", want: true, why: "the external test package"},
		{name: "main", path: "m/internal/app/commands/greet.test", want: true, why: "the test-main package"},
		{name: "greet", path: "m/internal/app/commands/greet", want: false, why: "the real package is judged"},
		{name: "greet", path: "m/internal/app/commands/greet_testing", want: false, why: "a suffix that merely resembles one"},
	} {
		assert.Equal(t, tc.want, isScaffoldingPackage(passWith(tc.name, tc.path)),
			"isScaffoldingPackage(%s, %s): %s", tc.name, tc.path, tc.why)
	}
}

// passWith builds the minimal pass isScaffoldingPackage inspects: a package
// with a name (the file's package clause) and an import path.
func passWith(name, path string) *analysis.Pass {
	return &analysis.Pass{Pkg: types.NewPackage(path, name)}
}

// TestRunJudgesOnlyRealCommandPackages names run's guard, which is three
// conditions doing three different jobs. Skipping the driver-synthesized
// scaffolding stops the analyzer reporting a missing entry point against a
// package that exists only in the driver; skipping non-command packages keeps
// per-command obligations off ordinary code; and the len(pass.Files) == 0 test
// is what lets checkCommandFunc index pass.Files without a bounds check — a
// test-only directory whose sole files are external tests yields an empty pass,
// and indexing it would panic inside the linter.
func TestRunJudgesOnlyRealCommandPackages(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		path string
		why  string
	}{
		{name: "greet_test", path: "m/internal/app/commands/greet", why: "an external test package carries no command source"},
		{name: "main", path: "m/internal/app/commands/greet.test", why: "the test-main package carries none either"},
		{name: "render", path: "m/internal/app/commands/greet/internal/render", why: "a nested helper is not a command"},
		{name: "widget", path: "m/pkg/widget", why: "an unrelated package"},
	} {
		pass := passWith(tc.name, tc.path)

		_, err := run(pass)

		assert.NoError(t, err, "%s: %s", tc.path, tc.why)
	}

	// A real command package with no files must also be skipped rather than
	// indexed — this is the bounds guard, and reaching it is the only way to
	// know it holds.
	assert.NotPanics(t, func() {
		_, _ = run(passWith("greet", "m/internal/app/commands/greet"))
	}, "an empty file list must be skipped, not indexed")
}
