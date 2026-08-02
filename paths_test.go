package cliapp

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

// TestIsCommandTreeMatchesAnyDepthBelowTheMarker names isCommandTree's claim.
// Every package beneath internal/app/commands is in scope AT ANY DEPTH — the
// depth-1 predicate this replaces made `app tenant create` invisible to every
// check. Depth is not the command/helper discriminator; the declaration gate in
// run is, so a nested helper passing this predicate is still exempt.
func TestIsCommandTreeMatchesAnyDepthBelowTheMarker(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path importPath
		why  string
		want bool
	}{
		{path: "m/internal/app/commands/greet", want: true, why: "one segment below the marker"},
		{path: "example.com/x/internal/app/commands/serve", want: true, why: "the module prefix is irrelevant"},
		{path: "m/internal/app/commands/tenant/create", want: true, why: "a nested command is in scope"},
		{path: "m/internal/app/commands/greet/internal/render", want: true, why: "a nested helper is in scope here; the declaration gate exempts it"},

		{path: "m/internal/app/commands", want: false, why: "the marker directory itself holds no command"},
		{path: "m/internal/app/commands/", want: false, why: "an empty segment names no command"},
		{path: "m/internal/app/command/greet", want: false, why: "the marker is 'commands', not 'command'"},
		{path: "m/internal/commands/greet", want: false, why: "the marker is the full app/commands path"},
		{path: "m/pkg/greet", want: false, why: "an unrelated package"},
		{path: "", want: false, why: "an empty path"},
	} {
		assert.Equal(t, tc.want, isCommandTree(tc.path), "isCommandTree(%q): %s", tc.path, tc.why)
	}
}

// TestIsScaffoldingPackageSkipsTheDriverSynthesizedPasses names the guard run
// depends on. The analysis driver synthesizes an external test package and a
// test-main package for any package that has tests, and NEITHER carries the
// command source — so each would be misjudged against a package that does not
// exist on disk.
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

// TestRunJudgesOnlyRealCommandPackages names run's guard: the
// driver-synthesized scaffolding is skipped, packages outside the command tree
// are skipped, and an empty file list is skipped rather than scanned — each
// without error and without a report.
func TestRunJudgesOnlyRealCommandPackages(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		path string
		why  string
	}{
		{name: "greet_test", path: "m/internal/app/commands/greet", why: "an external test package carries no command source"},
		{name: "main", path: "m/internal/app/commands/greet.test", why: "the test-main package carries none either"},
		{name: "widget", path: "m/pkg/widget", why: "an unrelated package"},
	} {
		pass := passWith(tc.name, tc.path)

		_, err := run(pass)

		assert.NoError(t, err, "%s: %s", tc.path, tc.why)
	}

	// A command-tree package with no files must also be skipped rather than
	// scanned — reaching the guard is the only way to know it holds.
	assert.NotPanics(t, func() {
		_, _ = run(passWith("greet", "m/internal/app/commands/greet"))
	}, "an empty file list must be skipped, not scanned")
}
