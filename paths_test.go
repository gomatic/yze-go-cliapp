package cliapp

// White-box tests for the classification rules. They decide WHETHER a package
// is judged at all and WHICH imports the alias rule judges, so a false positive
// imposes command-package obligations on a helper that has none, and a false
// negative exempts a real command from every check in this analyzer — silently,
// with a green gate.

import (
	"go/ast"
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
		{path: "m/myinternal/app/commands/greet", want: false, why: "the marker is a leading segment, not a substring"},
		{path: "m/pkg/greet", want: false, why: "an unrelated package"},
		{path: "", want: false, why: "an empty path"},
	} {
		assert.Equal(t, tc.want, isCommandTree(tc.path), "isCommandTree(%q): %s", tc.path, tc.why)
	}
}

// TestVerbIsThePackageName names the clause the package doc states and that
// nothing used to read: the verb IS the package name. Every judgement that
// depends on it — which entry point stays, which file is the command file, and
// what the diagnostic tells the author to move — is decided here rather than by
// filename order.
func TestVerbIsThePackageName(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		fn   string
		pkg  packageName
		verb commandVerb
		why  string
		want bool
	}{
		{fn: "Command", pkg: "greet", verb: "greet", want: true, why: "the bare spelling declares the package's own verb"},
		{fn: "GreetCommand", pkg: "greet", verb: "greet", want: true, why: "the explicit spelling of the same verb"},
		{fn: "GREETCommand", pkg: "greet", verb: "greet", want: true, why: "the verb is exported and the package name is not, so the compare folds case"},

		{fn: "PlanCommand", pkg: "mismatch", verb: "plan", want: false, why: "a verb the package is not named for"},
		{fn: "ApplyCommand", pkg: "config", verb: "apply", want: false, why: "the verb names the nested package it belongs in"},
		{fn: "Command", pkg: "config", verb: "config", want: true, why: "the resident of a parent command that mounts subcommands"},
		{fn: "PlanCommand", pkg: "Plan", verb: "plan", want: true, why: "the compare folds case in both directions: an unconventionally-capitalised package still names the verb, and the remedy this rule prescribes is a move, which is not the edit that package needs"},
	} {
		fn := &ast.FuncDecl{Name: ast.NewIdent(tc.fn)}

		assert.Equal(t, tc.verb, verbOf(fn, tc.pkg), "verbOf(%s, %s): %s", tc.fn, tc.pkg, tc.why)
		assert.Equal(t, tc.want, isPackageVerb(fn, tc.pkg), "isPackageVerb(%s, %s): %s", tc.fn, tc.pkg, tc.why)
	}
}

// TestKeywordVerbsAreTheOnesNoPackageCanHave names the exemption's matcher. The
// rule is "the verb is the package name", and for a verb Go reserves there is
// no source that satisfies it — `package import` does not parse — so an author
// whose CLI verb is one has already spelled the package as closely as Go
// permits and has no edit left to make. The exemption is keyed on token.Lookup
// and not on a name that reads keyword-ish, which is what the picker fixture
// pins from the other side.
func TestKeywordVerbsAreTheOnesNoPackageCanHave(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		verb commandVerb
		why  string
		want bool
	}{
		{verb: "import", want: true, why: "a plausible CLI verb Go reserves"},
		{verb: "select", want: true, why: "likewise"},
		{verb: "go", want: true, why: "likewise, and the shortest"},
		{verb: "type", want: true, why: "likewise"},
		{verb: "range", want: true, why: "likewise"},

		{verb: "importer", want: false, why: "a package name Go permits, and the honest spelling of the import verb"},
		{verb: "selector", want: false, why: "resembling a keyword is not being one"},
		{verb: "apply", want: false, why: "an ordinary verb"},
		{verb: "", want: false, why: "no verb at all"},
	} {
		assert.Equal(t, tc.want, isKeywordVerb(tc.verb), "isKeywordVerb(%q): %s", tc.verb, tc.why)
	}
}

// TestResidentCommandIsNotDecidedByFileOrder names which entry point stays when
// a package declares more than one. Reporting everything after the first in
// FILE order tells an author holding Command and ApplyCommand to move Command —
// into a package whose verb would be "command", which the same rule forbids. The
// resident is the one carrying the package's verb, and the bare spelling wins
// so that two candidates do not put the choice back on a filename.
func TestResidentCommandIsNotDecidedByFileOrder(t *testing.T) {
	t.Parallel()

	fn := func(name string) *ast.FuncDecl { return &ast.FuncDecl{Name: ast.NewIdent(name)} }

	apply, command, dup, plan := fn("ApplyCommand"), fn("Command"), fn("DupCommand"), fn("PlanCommand")

	for _, tc := range []struct {
		want     *ast.FuncDecl
		name     string
		pkg      packageName
		why      string
		commands []*ast.FuncDecl
	}{
		{name: "sole", pkg: "greet", commands: []*ast.FuncDecl{command}, want: command, why: "the only entry point carries the package verb"},
		{name: "late", pkg: "config", commands: []*ast.FuncDecl{apply, command}, want: command, why: "the resident is not the first in file order"},
		{name: "bare-wins", pkg: "dup", commands: []*ast.FuncDecl{dup, command}, want: command, why: "the bare spelling is preferred over an equally valid explicit one"},
		{name: "none", pkg: "multiverb", commands: []*ast.FuncDecl{apply, plan}, want: nil, why: "no entry point carries the package verb"},
	} {
		assert.Same(t, tc.want, residentCommand(tc.commands, tc.pkg), "%s: %s", tc.name, tc.why)
	}
}

// passWith builds the minimal pass run inspects before it reads any file: a
// package with a name (the file's package clause) and an import path.
func passWith(name, path string) *analysis.Pass {
	return &analysis.Pass{Pkg: types.NewPackage(path, name)}
}

// TestRunJudgesOnlyRealCommandPackages names run's two gates. A package outside
// the command tree is skipped, and a package inside it that declares no entry
// point is skipped by the declaration gate rather than by a file count — which
// is the gate that has to hold, because it is the one standing between an
// empty package and a nil command file. Reaching it is the only way to know it
// does.
func TestRunJudgesOnlyRealCommandPackages(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		path string
		why  string
	}{
		{name: "widget", path: "m/pkg/widget", why: "an unrelated package"},
		{name: "lookalike", path: "m/myinternal/app/commands/lookalike", why: "a path containing the marker without carrying it as a segment"},
		{name: "greet_test", path: "m/internal/app/commands/greet", why: "the external test package carries no command source"},
		{name: "main", path: "m/internal/app/commands/greet.test", why: "the test-main package carries none either"},
		{name: "greet", path: "m/internal/app/commands/greet", why: "a command-tree package with no files declares no entry point"},
	} {
		assert.NotPanics(t, func() {
			_, err := run(passWith(tc.name, tc.path))

			assert.NoError(t, err, "%s: %s", tc.path, tc.why)
		}, "%s: %s", tc.path, tc.why)
	}
}
