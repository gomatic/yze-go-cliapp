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
	messageAlias = "command package: import the domain package with the %q alias"
	messageTaken = "command package: %q already binds %q here, so the domain package cannot take it — rename that import first"
)

// domainPaths are the two import paths a command package may name as its own
// domain package, both derived from the ANALYZED package's path — which is what
// the build compiled it as, and which no judged file can rewrite.
//
// The counterpart is the domain package that corresponds to this command,
// segment for segment: internal/app/commands/tenant/create is backed by
// internal/domain/tenant/create. The tier root is the fallback, and it is only
// the domain package of a module laying its whole domain out as one package at
// internal/domain — which the reference layout permits and which a file
// importing no counterpart is presumed to be.
type domainPaths struct {
	counterpart importPath
	tierRoot    importPath
}

// domainPathsOf derives the two candidate paths from a command package's own
// path. It is only called after isCommandTree, so the marker is present and the
// remainder names a command.
//
// Deriving them by exact path — rather than matching every import at or beneath
// the tier — is what bounds the rule to one import per file. It also retires a
// whole class of boundary mistake for free: internal/domainhelpers and
// internal/domains cannot equal either path, so no segment-matching predicate
// is needed to keep them out, and there is no substring test left to get wrong
// in either direction.
func domainPathsOf(pkgPath importPath) domainPaths {
	namespace, command, _ := strings.Cut(string(pkgPath), commandsMarker)
	return domainPaths{
		counterpart: importPath(namespace + domainMarker + "/" + command),
		tierRoot:    importPath(namespace + domainMarker),
	}
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
	paths := domainPathsOf(importPath(pass.Pkg.Path()))
	for _, file := range productionFiles(pass) {
		checkFileDomainAliases(pass, file, paths)
	}
}

// checkFileDomainAliases judges the ONE import a file can be asked to bind
// "domain" to: the package's own domain counterpart, or — when the file names
// no counterpart — the tier root, which is the whole domain of a module that
// has only one.
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
func checkFileDomainAliases(pass *analysis.Pass, file *ast.File, paths domainPaths) {
	specs := entitledImports(file, paths)
	named := firstNamedImport(specs)
	if named == nil || anyBindsDomain(pass, specs) {
		return
	}
	if holder := domainHolder(pass, file); holder != nil {
		pass.Reportf(named.Pos(), messageTaken, importedPath(holder), domainAlias)
		return
	}
	pass.Reportf(named.Pos(), messageAlias, domainAlias)
}

// entitledImports returns the specs naming the package's own domain package:
// its counterpart where the file imports it, and the tier root only where it
// does not. The counterpart wins because it is the package the standard names;
// the root is the domain package of a flat module and nobody's counterpart in a
// grouped one, so letting it stand in while a counterpart is present would sell
// silence for one import line.
func entitledImports(file *ast.File, paths domainPaths) []*ast.ImportSpec {
	if specs := importsOf(file, paths.counterpart); len(specs) > 0 {
		return specs
	}
	return importsOf(file, paths.tierRoot)
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

// firstNamedImport returns the first spec that binds an identifier, or nil when
// every spec is blank — the shape a command mounting registration side effects
// writes, and the one where the remedy has nothing to rename.
func firstNamedImport(specs []*ast.ImportSpec) *ast.ImportSpec {
	for _, imp := range specs {
		if !isBlankImport(imp) {
			return imp
		}
	}
	return nil
}

// anyBindsDomain reports whether the domain package is bound to the standard's
// name anywhere in the file. One binding satisfies the rule's reason — the call
// sites reading domain.X have it — so a redundant second import of the same
// path under another name is a wart for some other rule, not an alias
// violation, and reporting it would prescribe a redeclaration.
func anyBindsDomain(pass *analysis.Pass, specs []*ast.ImportSpec) bool {
	for _, imp := range specs {
		if bindsDomain(pass, imp) {
			return true
		}
	}
	return false
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
