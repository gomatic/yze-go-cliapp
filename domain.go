package cliapp

// The domain import rule. Which import it judges is a scope decision, and the
// two exemptions below are not judgements about intent at all — each is a place
// where the remedy the diagnostic prescribes would not compile, which is the
// one thing an author cannot answer except with a disablement.

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// domainMarker is the path of the domain tier below the namespace a command
// tree sits in.
const domainMarker = "/internal/domain"

// isDomainImport reports whether an import path names the domain tier itself or
// a package beneath it. The doc comment says "the domain package" and names no
// discriminator, so this is the one it means, and its two boundaries are stated
// here because a mistake at either is invisible in a diff of findings.
//
// The match is on whole path SEGMENTS. A plain substring test is wrong in both
// directions at once: requiring a trailing slash exempts a module whose whole
// domain is one package at internal/domain — which the reference layout has,
// holding the vocabulary every Run contract shares — and dropping it admits a
// sibling internal/domainhelpers, whose alias nothing has any business
// prescribing.
//
// Nothing further narrows it to the analyzed module's own tier. Deriving the
// tier from the command package's own path reads as tighter and is not: it
// silently exempts a command tree that is not at the module root, which is the
// scope mistake this file exists to avoid. The narrowing would also be
// unreachable: the go command's internal-visibility rule restricts a path
// holding an internal segment to importers within the namespace above it, so a
// module conjunct here has no input that reaches it, and a guard nothing can
// reach is dead code rather than defence.
func isDomainImport(path importPath) bool {
	_, rest, found := strings.Cut(string(path), domainMarker)
	return found && (rest == "" || strings.HasPrefix(rest, "/"))
}

// isDomainTierRoot reports whether an import path names the tier itself rather
// than a group beneath it. The root is the shared vocabulary package and so is
// nobody's domain counterpart, which is why the collision exemption refuses to
// let it stand in for one. TestDomainImportIsMatchedOnSegments pins the
// boundary it reads.
func isDomainTierRoot(path importPath) bool {
	_, rest, found := strings.Cut(string(path), domainMarker)
	return found && rest == ""
}

// importedPath is the path an import spec names. ImportSpec.Path is a quoted
// STRING LITERAL, and Go admits the raw form as readily as the interpreted one,
// so the quotes come off through strconv.Unquote rather than by trimming the
// double quote — trimming leaves a backquoted path carrying its backquotes and
// silently outside every match here, which is a domain import turned off by one
// keystroke. A literal Unquote rejects is one the go parser rejected first, so
// it reaches no compiled package; it yields the empty path, which no predicate
// here matches. TestImportedPathReadsBothStringForms pins both forms.
func importedPath(imp *ast.ImportSpec) importPath {
	path, err := strconv.Unquote(imp.Path.Value)
	if err != nil {
		return ""
	}
	return importPath(path)
}

// checkDomainAlias reports domain imports not aliased "domain" — in PRODUCTION
// files only. A test file may alias the domain import to disambiguate
// (greetdomain) because it names packages the command file does not; the shape
// this rule judges is the command package's, and a test file carries none of
// it. That is a different reason from the one that keeps test files out of the
// entry-point scan, which is misclassification; the two skips share the
// isTestFile filter and nothing else.
func checkDomainAlias(pass *analysis.Pass) {
	for _, file := range productionFiles(pass) {
		checkFileDomainAliases(pass, file)
	}
}

// checkFileDomainAliases judges one file's domain-tier imports. "domain" is a
// FILE-scoped name, so telling an import to take it when another import in the
// same file already has it prescribes a redeclaration, and the file stops
// compiling — measured by building the remedy, not by reading the spec. The shape that exemption exists for is one file naming a
// domain group AND a type package nested beneath it, so that is what it is
// keyed on: the holder of the alias is a STRICT PATH PREFIX of the import it
// excuses, and is not the tier root.
//
// Both narrowings were bought by a forgery an earlier draft admitted. Requiring
// only "some domain-tier import binds domain" let a file silence its own
// misaliased group by importing the tier root as "domain" — the root is nobody's
// counterpart — or by importing the SAME path twice, once as "domain" and once
// under the wrong name, leaving every call site untouched. Two lines, total
// silence, and the property acquired was none.
//
// What remains forgeable is stated rather than claimed away: for a group nested
// under another group, importing the PARENT group as "domain" and using it does
// excuse the child. It costs a real import of a real ancestor domain package,
// which is closer to the property than to a marker for it, and no shorter than
// aliasing the child correctly.
func checkFileDomainAliases(pass *analysis.Pass, file *ast.File) {
	imports := domainImports(file)
	for _, imp := range imports {
		if isBlankImport(imp) || aliasedDomain(imp) || excusedByAncestor(imports, imp) {
			continue
		}
		pass.Reportf(imp.Pos(), "command package: import the domain package with the \"domain\" alias")
	}
}

// excusedByAncestor reports whether another import in the same file holds the
// "domain" alias for a group this one is nested beneath.
func excusedByAncestor(imports []*ast.ImportSpec, imp *ast.ImportSpec) bool {
	path := importedPath(imp)
	for _, other := range imports {
		holder := importedPath(other)
		if aliasedDomain(other) && !isDomainTierRoot(holder) && strings.HasPrefix(string(path), string(holder)+"/") {
			return true
		}
	}
	return false
}

// domainImports collects one file's imports at or beneath the domain tier.
func domainImports(file *ast.File) []*ast.ImportSpec {
	var found []*ast.ImportSpec
	for _, imp := range file.Imports {
		if isDomainImport(importedPath(imp)) {
			found = append(found, imp)
		}
	}
	return found
}

// aliasedDomain reports whether an import binds the name "domain".
func aliasedDomain(imp *ast.ImportSpec) bool {
	return imp.Name != nil && imp.Name.Name == "domain"
}

// isBlankImport reports whether an import binds no identifier at all, which a
// command mounting a domain's registration side effects writes. There is no
// name in it to spell "domain" with: aliasing it leaves an import nothing uses,
// and the file stops compiling.
func isBlankImport(imp *ast.ImportSpec) bool {
	return imp.Name != nil && imp.Name.Name == "_"
}
