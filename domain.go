package cliapp

// The domain import rule. Which import it judges is a scope decision, and it is
// the whole of the difficulty: "domain" is a FILE-SCOPED name and a file has
// exactly one, so a rule that asks more than one import for it prescribes a
// redeclaration — an instruction whose only answer is a disablement. Everything
// here exists to make the rule ask exactly one import per file, and to make the
// remedy compile when the name is already spoken for.

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// domainMarker is the path of the domain tier below the namespace a command
// tree sits in.
const domainMarker = "/internal/domain"

// domainAlias is the name the standard gives a command package's own domain
// package, so that every call site in the app tier reads domain.Run whatever
// verb the package answers to.
const domainAlias = "domain"

// blankName is the import name that is not a name: a blank import binds no
// identifier to spell the alias with. A dot import is not in this company — it
// binds every exported name of the package instead of one, which is an edit the
// author can still make, so it is judged like any other.
const blankName = "_"

// Diagnostic messages. The second exists because the first is not always an
// edit that compiles: where another import in the file already binds "domain",
// aliasing this one redeclares it, and the author's only remaining moves would
// be a baseline or a //nolint.
const (
	messageAlias     = "command package: import the domain package with the %q alias"
	messageTaken     = "command package: %q already binds %q here, so the domain package cannot take it — rename that import first"
	messageDuplicate = "command package: the domain package is already imported as %q in this file — drop this second import and use it"
)

// counterpartOf derives the import path of a command package's OWN domain
// package from the ANALYZED package's path — which is what the build compiled
// it as, and which no judged file can rewrite. The counterpart corresponds
// segment for segment: internal/app/commands/tenant/create is backed by
// internal/domain/tenant/create. It is only called after isCommandTree, so the
// marker is present and the remainder names a command.
//
// The FIRST marker decides, because the tier a path belongs to has to be the
// one yze/clidomain resolves the shared vocabulary against, and its
// sharedVocabulary cuts at the first occurrence too (yze-go-clidomain
// spelling.go:98). Two analyzers disagreeing about where the tier is, is how a
// file-scoped name ends up demanded twice.
//
// Deriving one exact path — rather than matching every import at or beneath the
// tier — is what bounds the rule to one import per file. It also retires a
// whole class of boundary mistake for free: internal/domainhelpers and
// internal/domains cannot equal it, so no segment-matching predicate is needed
// to keep them out, and there is no substring test left to get wrong in either
// direction.
//
// Nothing stands in for a counterpart that is not imported. The shared
// vocabulary package at the tier root is the obvious candidate and is wrong:
// it is every tier's vocabulary and nobody's counterpart, so a mount-only
// parent command naming domain.Argument, or the SECOND FILE of a command
// package whose counterpart is imported in the first, would both be told to
// give it the name — and taking that instruction makes "domain" mean two
// different packages inside one package, which is the collision this rule was
// rewritten to stop prescribing.
func counterpartOf(pkgPath importPath) importPath {
	namespace, command, _ := strings.Cut(string(pkgPath), commandsMarker)
	return importPath(namespace + domainMarker + "/" + command)
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

// boundName is the identifier an import actually binds in its file — which is
// what the rule is about, and what an alias is only one way to set. An import
// with no alias binds the imported package's OWN name, so an unaliased import
// of the tier root binds "domain" already, and every call site in the file
// reads domain.X exactly as the standard asks. Reading the alias instead
// reports that file for an edit whose whole content is a redundant alias.
//
// The name comes from the imported package rather than from the last path
// segment, because those differ and only the first is what the file binds.
// Forging it means declaring `package domain`, which acquires the property the
// rule exists for rather than a marker standing for it.
func boundName(pass *analysis.Pass, imp *ast.ImportSpec) string {
	if imp.Name != nil {
		return imp.Name.Name
	}
	path := importedPath(imp)
	for _, imported := range pass.Pkg.Imports() {
		if importPath(imported.Path()) == path {
			return imported.Name()
		}
	}
	return ""
}

// bindsDomain reports whether an import binds the standard's name.
func bindsDomain(pass *analysis.Pass, imp *ast.ImportSpec) bool {
	return boundName(pass, imp) == domainAlias
}

// isBlankImport reports whether an import binds no identifier at all, which a
// command mounting a domain's registration side effects writes. There is no
// name in it to spell "domain" with: aliasing it leaves an import nothing uses,
// and the file stops compiling.
func isBlankImport(imp *ast.ImportSpec) bool {
	return imp.Name != nil && imp.Name.Name == blankName
}

// checkDomainAlias reports domain imports not bound to "domain" — in PRODUCTION
// files only. A test file may alias the domain import to disambiguate
// (greetdomain) because it names packages the command file does not; the shape
// this rule judges is the command package's, and a test file carries none of
// it. That is a different reason from the one that keeps test files out of the
// entry-point scan, which is misclassification; the two skips share the
// isTestFile filter and nothing else.
func checkDomainAlias(pass *analysis.Pass) {
	counterpart := counterpartOf(importPath(pass.Pkg.Path()))
	for _, file := range productionFiles(pass) {
		checkFileDomainAliases(pass, file, counterpart)
	}
}

// checkFileDomainAliases judges the ONE import a file can be asked to bind
// "domain" to: the package's own domain counterpart, and nothing else ever.
//
// Every other import at or beneath the tier is out of the rule's reach, and
// that is the point rather than a concession. A command file legitimately names
// the shared vocabulary package alongside its counterpart, or a type package
// nested under its own group; there is one "domain" to give and the counterpart
// has it, so asking a second import for the same name is asking for source that
// does not compile. Whether a command package should be reaching into a domain
// group that is not its own is cross-package correspondence, which this
// analyzer's doc comment hands to stickler/clilayout.
//
// The finding sits on the FIRST spec that binds a name, so that applying the
// remedy ends the finding: reporting every spec of the path would leave a
// second report standing after the first was answered, and the second answer
// does not exist.
func checkFileDomainAliases(pass *analysis.Pass, file *ast.File, counterpart importPath) {
	named := namedImports(importsOf(file, counterpart))
	if len(named) == 0 {
		return
	}
	if bound := boundImport(pass, named); bound != nil {
		reportDuplicates(pass, named, bound)
		return
	}
	if holder := domainHolder(pass, file); holder != nil {
		pass.Reportf(named[0].Pos(), messageTaken, importedPath(holder), domainAlias)
		return
	}
	pass.Reportf(named[0].Pos(), messageAlias, domainAlias)
}

// importsOf collects the specs in a file naming one exact import path. Go
// permits the same path twice under two names, so this is a list.
func importsOf(file *ast.File, path importPath) []*ast.ImportSpec {
	var found []*ast.ImportSpec
	for _, imp := range file.Imports {
		if importedPath(imp) == path {
			found = append(found, imp)
		}
	}
	return found
}

// namedImports drops the specs that bind no identifier. A blank import has no
// name in it to spell the alias with, and a file whose every spec is blank —
// the shape a command mounting registration side effects writes — is asked for
// nothing.
func namedImports(specs []*ast.ImportSpec) []*ast.ImportSpec {
	var found []*ast.ImportSpec
	for _, imp := range specs {
		if !isBlankImport(imp) {
			found = append(found, imp)
		}
	}
	return found
}

// boundImport returns the spec that binds the standard's name, searching every
// spec of the path rather than the first: Go permits the same path twice under
// two names, and the order of an import block is gofmt's rather than the
// author's.
func boundImport(pass *analysis.Pass, specs []*ast.ImportSpec) *ast.ImportSpec {
	for _, imp := range specs {
		if bindsDomain(pass, imp) {
			return imp
		}
	}
	return nil
}

// reportDuplicates reports every OTHER spec of the domain package's path once
// one of them binds the name. Without this, two lines buy silence: import your
// own domain package a second time as "domain", use it once, and every call
// site may keep reading greetdomain. The name is bound, so the alias remedy
// would redeclare it — the takeable edit is to drop the second import, and the
// message says so rather than prescribing source that does not build.
func reportDuplicates(pass *analysis.Pass, named []*ast.ImportSpec, bound *ast.ImportSpec) {
	for _, imp := range named {
		if imp != bound {
			pass.Reportf(imp.Pos(), messageDuplicate, domainAlias)
		}
	}
}

// domainHolder returns the import already binding "domain" in this file, if
// any. It is consulted only once the domain package is known not to be that
// import, so what it finds is a squatter — the shared vocabulary package under
// an unaliased import, or an unrelated package aliased "domain" — and naming it
// is what turns an impossible instruction into two edits that both compile.
func domainHolder(pass *analysis.Pass, file *ast.File) *ast.ImportSpec {
	for _, imp := range file.Imports {
		if bindsDomain(pass, imp) {
			return imp
		}
	}
	return nil
}
