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

// TestDomainImportIsMatchedOnSegments names the boundary the alias rule turns
// on. The doc comment says "the domain package" and names no discriminator, so
// this is the one it means — and a substring test gets it wrong in both
// directions at once. Requiring a trailing slash exempts a module whose whole
// domain is one package at internal/domain, silently, from the only rule this
// analyzer has for it; dropping it admits internal/domainhelpers, whose alias
// nothing has any business prescribing. Neither mistake reports anything, so
// neither shows up in a diff of findings.
func TestDomainImportIsMatchedOnSegments(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path importPath
		why  string
		want bool
	}{
		{path: "example.com/x/internal/domain", want: true, why: "the tier root is a real package in the reference layout"},
		{path: "example.com/x/internal/domain/greet", want: true, why: "a domain group"},
		{path: "example.com/x/internal/domain/config/get", want: true, why: "a nested domain group"},
		{path: "example.com/x/internal/domain/greet/model", want: true, why: "a type package beneath a group is still in the tier"},

		{path: "example.com/x/internal/domainhelpers", want: false, why: "a sibling whose name merely begins with the tier's"},
		{path: "example.com/x/internal/domains", want: false, why: "the segment is 'domain', not a prefix of one"},
		{path: "example.com/x/internal/app/commands/tenant", want: false, why: "the command tree is not the domain tier"},
		{path: "example.com/x/domain", want: false, why: "a domain package outside any internal tier"},
		{path: "", want: false, why: "an empty path"},
	} {
		assert.Equal(t, tc.want, isDomainImport(tc.path), "isDomainImport(%q): %s", tc.path, tc.why)
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
		verb string
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
