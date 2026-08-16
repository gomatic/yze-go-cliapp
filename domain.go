package cliapp

// The domain import rule. Which import it judges is a scope decision, and the
// two exemptions below are not judgements about intent at all — each is a place
// where the remedy the diagnostic prescribes would not compile, which is the
// one thing an author cannot answer except with a disablement.

import (
	"go/ast"
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
// FILE-scoped name, so the judgement is per file: once one import in the file
// binds it, no other import can, and asking a second one to take it prescribes
// a redeclaration. That exemption is not a judgement about which package
// deserves the alias — forging it means importing a real domain package under
// the alias the rule wants, which is the property, not a marker for it.
func checkFileDomainAliases(pass *analysis.Pass, file *ast.File) {
	imports := domainImports(file)
	for _, imp := range imports {
		if aliasedDomain(imp) {
			return
		}
	}
	for _, imp := range imports {
		if isBlankImport(imp) {
			continue
		}
		pass.Reportf(imp.Pos(), "command package: import the domain package with the \"domain\" alias")
	}
}

// domainImports collects one file's imports at or beneath the domain tier.
func domainImports(file *ast.File) []*ast.ImportSpec {
	var found []*ast.ImportSpec
	for _, imp := range file.Imports {
		if isDomainImport(importPath(strings.Trim(imp.Path.Value, `"`))) {
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
