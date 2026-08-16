package cliapp

// White-box tests for the domain-import rule: which ONE import in a file may be
// asked to bind "domain", what an import actually binds, and what the two
// diagnostics say. The corpus cannot see a message and cannot vary an imported
// package's own package clause, so both are owned here.

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

// TestDomainPathsOfDerivesBothCandidatesFromTheCommandsOwnPath names the scope decision the
// whole domain rule turns on. "domain" is a file-scoped name and a file has
// exactly one, so at most one import per file can be asked for it; which one is
// decided HERE, from the analyzed package's own import path, which is what the
// build compiled it as and which no judged file can rewrite.
//
// Deriving the counterpart segment for segment is also what retires the
// substring-versus-segment boundary the old scope predicate had to defend in
// both directions at once: internal/domainhelpers and internal/domains cannot
// equal either derived path, so there is no match left to widen.
func TestDomainPathsOfDerivesBothCandidatesFromTheCommandsOwnPath(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		pkg         importPath
		counterpart importPath
		tierRoot    importPath
		why         string
	}{
		{
			pkg: "m/internal/app/commands/greet", counterpart: "m/internal/domain/greet", tierRoot: "m/internal/domain",
			why: "a top-level command is backed by the domain package of the same name",
		},
		{
			pkg: "example.com/x/internal/app/commands/tenant/create", counterpart: "example.com/x/internal/domain/tenant/create", tierRoot: "example.com/x/internal/domain",
			why: "a nested command is backed segment for segment, at the same depth",
		},
		{
			pkg: "m/internal/app/commands/greet/internal/render", counterpart: "m/internal/domain/greet/internal/render", tierRoot: "m/internal/domain",
			why: "a helper beneath a command derives paths it never imports; the declaration gate is what exempts it, not this",
		},
	} {
		paths := domainPathsOf(tc.pkg)

		assert.Equal(t, tc.counterpart, paths.counterpart, "counterpart of %q: %s", tc.pkg, tc.why)
		assert.Equal(t, tc.tierRoot, paths.tierRoot, "tier root of %q: %s", tc.pkg, tc.why)
		assert.NotEqual(
			t,
			importPath(string(tc.tierRoot)+"helpers"),
			paths.counterpart,
			"a sibling of the tier is never the counterpart",
		)
		assert.NotEqual(
			t,
			importPath(string(tc.tierRoot)+"helpers"),
			paths.tierRoot,
			"a sibling of the tier is never the tier root",
		)
	}
}

// TestImportedPathReadsBothStringForms names what an import path IS. Go admits
// the raw string form as readily as the interpreted one and the build treats
// them identically, so trimming the double quote leaves a backquoted path
// carrying its backquotes, outside every predicate here — a domain import
// turned off by one keystroke. This is asserted against a constructed
// ImportSpec rather than only through the rawstring fixture, because gofmt
// rewrites a raw import path back to the interpreted form: the fixture's
// discriminating property is one a routine format deletes, and a case with a
// half-life is not a case.
func TestImportedPathReadsBothStringForms(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		literal string
		want    importPath
		why     string
	}{
		{literal: `"m/internal/domain/greet"`, want: "m/internal/domain/greet", why: "the interpreted form"},
		{literal: "`m/internal/domain/greet`", want: "m/internal/domain/greet", why: "the raw form the build accepts identically"},
		{literal: "`m/internal/domain`", want: "m/internal/domain", why: "the raw form of the tier root, where trimming fails outright"},
		{literal: `"m/internal/domain"`, want: "m/internal/domain", why: "the interpreted tier root"},
		{literal: "m/unquoted", want: "", why: "a literal the go parser would have rejected first names no path"},
	} {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: tc.literal}}

		assert.Equal(t, tc.want, importedPath(spec), "importedPath(%s): %s", tc.literal, tc.why)
	}
}

// TestBoundNameIsTheIdentifierNotTheAlias names what the domain rule is
// actually about. An import with no alias binds the imported package's OWN
// name, so an unaliased import of a package whose clause is `package domain`
// already reads domain.X at every call site — the property the standard asks
// for, spelled without an alias. Reading ImportSpec.Name instead reports that
// file for an edit whose whole content is a redundant alias, and the two
// directions are held apart here because no fixture can vary an imported
// package's own clause without a second package to import.
func TestBoundNameIsTheIdentifierNotTheAlias(t *testing.T) {
	t.Parallel()

	vocabulary := types.NewPackage("m/internal/domain", "domain")
	counterpart := types.NewPackage("m/internal/domain/greet", "greet")

	pkg := types.NewPackage("m/internal/app/commands/greet", "greet")
	pkg.SetImports([]*types.Package{vocabulary, counterpart})
	pass := &analysis.Pass{Pkg: pkg}

	for _, tc := range []struct {
		alias string
		path  string
		want  string
		why   string
	}{
		{path: "m/internal/domain", want: "domain", why: "an unaliased import binds the imported package's own name"},
		{path: "m/internal/domain/greet", want: "greet", why: "the same rule, where that name is not domain"},
		{alias: "domain", path: "m/internal/domain/greet", want: "domain", why: "an explicit alias binds what it spells"},
		{alias: "greetdomain", path: "m/internal/domain/greet", want: "greetdomain", why: "likewise when it spells something else"},
		{alias: "_", path: "m/internal/domain/greet", want: "_", why: "a blank import binds no identifier, and the blank is what says so"},
		{alias: ".", path: "m/internal/domain/greet", want: ".", why: "a dot import binds every exported name rather than one"},
		{path: "m/internal/domain/absent", want: "", why: "a path the package does not import names nothing this can read"},
	} {
		imp := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(tc.path)}}
		if tc.alias != "" {
			imp.Name = ast.NewIdent(tc.alias)
		}

		assert.Equal(t, tc.want, boundName(pass, imp), "boundName(%q as %q): %s", tc.path, tc.alias, tc.why)
		assert.Equal(
			t,
			tc.want == domainAlias,
			bindsDomain(pass, imp),
			"bindsDomain(%q as %q): %s",
			tc.path,
			tc.alias,
			tc.why,
		)
	}
}

// TestMessageAliasAndMessageTakenEachNameAnEditThatCompiles owns the two
// diagnostic messages, which the corpus cannot see — its Finding is
// {rule, path, line} and carries no message at all. They are asserted from the
// standard rather than from a run, and the pair exists because one sentence
// cannot serve both arrangements: where "domain" is free, adding the alias is
// the edit; where another import already binds it, adding the alias is
// `domain redeclared in this block`, so the message names the holder and the
// edit is to rename that import first.
func TestMessageAliasAndMessageTakenEachNameAnEditThatCompiles(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		`command package: import the domain package with the "domain" alias`,
		fmt.Sprintf(messageAlias, domainAlias),
		"the free-name message asks for the alias and nothing else")

	assert.Equal(t,
		`command package: "m/internal/domain" already binds "domain" here, `+
			`so the domain package cannot take it — rename that import first`,
		fmt.Sprintf(messageTaken, importPath("m/internal/domain"), domainAlias),
		"the taken-name message names the holder, which is the import the author must edit")
}

// TestAnyBindsDomainReadsEverySpecOfThePath names why the search is over every
// spec of the domain package's path and not over the first. Go permits the same
// path twice under two names, and where one of them IS "domain" the rule's
// reason holds: the domain package is bound to the name. Reporting the other
// spec would prescribe `domain redeclared in this block`, and the only edit
// that answers it is a deletion no message here names.
//
// It is asserted against constructed specs rather than through a fixture
// because gofmt SORTS an import block, so a fixture cannot hold the
// domain-binding spec second — the discriminating property has a half-life of
// one format run, and a case with a half-life is not a case.
func TestAnyBindsDomainReadsEverySpecOfThePath(t *testing.T) {
	t.Parallel()

	pass := &analysis.Pass{Pkg: types.NewPackage("m/internal/app/commands/greet", "greet")}
	spec := func(alias string) *ast.ImportSpec {
		return &ast.ImportSpec{
			Name: ast.NewIdent(alias),
			Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote("m/internal/domain/greet")},
		}
	}

	for _, tc := range []struct {
		why   string
		specs []*ast.ImportSpec
		want  bool
	}{
		{specs: []*ast.ImportSpec{spec(domainAlias)}, want: true, why: "the only spec binds the name"},
		{specs: []*ast.ImportSpec{spec(domainAlias), spec("greetdomain")}, want: true, why: "the binding spec is first"},
		{specs: []*ast.ImportSpec{spec("greetdomain"), spec(domainAlias)}, want: true, why: "and the name is still bound when it is second"},
		{specs: []*ast.ImportSpec{spec("greetdomain")}, want: false, why: "no spec binds it"},
		{specs: nil, want: false, why: "the path is not imported at all"},
	} {
		assert.Equal(t, tc.want, anyBindsDomain(pass, tc.specs), "anyBindsDomain: %s", tc.why)
	}
}
