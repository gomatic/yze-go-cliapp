// Package cliapp provides a go/analysis analyzer enforcing the app tier of the
// opinionated three-tier CLI layout, for every command package beneath
// internal/app/commands/ at any depth. A command package declares itself by
// exporting a command entry point; for each such package the command file (the
// first one defining the entry point) leads with a const block, the package
// declares exactly ONE entry point — the verb is the package name, so a second
// verb belongs in its own nested package, as the reference layout demonstrates
// with config/{get,list,set} — and the domain package is imported under the
// "domain" alias.
//
// A package beneath the tree that declares no entry point — a helper nested
// under a command, or a shared grouping package — carries none of these
// obligations. Whether every self-declaring command has its domain counterpart
// (and vice versa) is cross-package correspondence, which is
// stickler/clilayout's job; the domain package's own contract is
// yze/clidomain's.
package cliapp

import (
	"go/ast"
	"go/token"
	"strings"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
)

// Analyzer reports per-package violations of the three-tier command-package layout.
var Analyzer = &analysis.Analyzer{
	Name: "cliapp",
	Doc:  "reports command packages that violate the app tier of the opinionated three-tier CLI layout",
	Run:  run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "cliapp",
	Categories: []goyze.Category{"cli", "structure"},
	URL:        "https://docs.gomatic.dev/yze/cliapp",
	Analyzer:   Analyzer,
}

// run checks self-declaring command packages for the per-package standards.
//
// The checks concern a command package's source code, so run skips the two
// scaffolding passes the driver synthesizes for a package that has tests — the
// external test package (clause "<pkg>_test") and the test-main package (import
// path "<pkg>.test") — neither of which carries the command source. The
// remaining pass(es) hold the command source; the driver collapses their
// identical diagnostics, so the checks run once.
//
// A package that declares no entry point is out of scope entirely: it is a
// helper or grouping package, and depth never makes a package a command. The
// converse gap — a command that SHOULD exist but declares nothing — is not a
// single-package fact and is stickler/clilayout's to report.
func run(pass *analysis.Pass) (any, error) {
	if isScaffoldingPackage(pass) || !isCommandTree(importPath(pass.Pkg.Path())) || len(pass.Files) == 0 {
		return nil, nil
	}
	commands := collectCommands(pass)
	if len(commands.funcs) == 0 {
		return nil, nil
	}
	reportNonConstFirst(pass, commands.file)
	checkSingleVerb(pass, commands.funcs)
	checkDomainAlias(pass)
	return nil, nil
}

// reportNonConstFirst reports when the command file's first non-import
// declaration is not a const block. The command file (the first one defining a
// command entry point) is the canonical metadata file, so the check targets it
// rather than an arbitrary first file of a multi-file package.
func reportNonConstFirst(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		if isImportDecl(decl) {
			continue
		}
		if !isConstDecl(decl) {
			pass.Reportf(decl.Pos(), "command package: the first declaration must be the const block")
		}
		return
	}
}

// isTestFile reports whether file is a test file (_test.go).
func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go")
}

func isImportDecl(decl ast.Decl) bool {
	gen, ok := decl.(*ast.GenDecl)
	return ok && gen.Tok == token.IMPORT
}

func isConstDecl(decl ast.Decl) bool {
	gen, ok := decl.(*ast.GenDecl)
	return ok && gen.Tok == token.CONST
}

// commandDecls is what a package declares of the command shape: every entry
// point in file order, and the command file — the first non-test file
// declaring one, nil only when funcs is empty.
type commandDecls struct {
	file  *ast.File
	funcs []*ast.FuncDecl
}

// collectCommands gathers the package's command declarations from its non-test
// files. Test files never carry the command source and are skipped: the
// in-package test-variant pass appends _test.go files to pass.Files, where an
// exported Test*Command test function (e.g. TestCommand) would otherwise be
// mistaken for an entry point — anchoring a spurious const-first diagnostic in
// the test file and marking a helper package as self-declaring. The result
// doubles as the self-declaration gate: no entry points marks a helper or
// grouping package.
func collectCommands(pass *analysis.Pass) commandDecls {
	var found commandDecls
	for _, file := range pass.Files {
		if isTestFile(pass, file) {
			continue
		}
		funcs := fileCommandFuncs(file)
		if len(funcs) > 0 && found.file == nil {
			found.file = file
		}
		found.funcs = append(found.funcs, funcs...)
	}
	return found
}

// fileCommandFuncs collects the command entry points one file declares.
func fileCommandFuncs(file *ast.File) []*ast.FuncDecl {
	var found []*ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := commandFunc(decl); ok {
			found = append(found, fn)
		}
	}
	return found
}

// commandFunc reports whether decl is a command entry point: an exported
// top-level function named Command or a <Verb>Command constructor. Recognition
// is deliberately broad — a package that exports PlanCommand has set out to be
// a command package and must be judged, not skipped — but declaring more than
// one is itself a violation; see checkSingleVerb.
func commandFunc(decl ast.Decl) (*ast.FuncDecl, bool) {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}
	return fn, fn.Recv == nil && ast.IsExported(fn.Name.Name) && strings.HasSuffix(fn.Name.Name, "Command")
}

// checkSingleVerb reports every command entry point beyond the first. One
// package is one verb: the reference layout nests each verb in its own package
// (config/get, config/list, config/set), so a package exporting
// PlanCommand/ApplyCommand holds verbs that each belong in a nested package of
// their own.
func checkSingleVerb(pass *analysis.Pass, commands []*ast.FuncDecl) {
	for _, fn := range commands[1:] {
		pass.Reportf(
			fn.Pos(),
			"command package: one verb per package — move %s into its own nested package",
			fn.Name.Name,
		)
	}
}

// checkDomainAlias reports domain imports not aliased as "domain" — in
// PRODUCTION files only. A test file may alias the domain import to
// disambiguate (greetdomain) without violating the command package's shape,
// exactly as test files are excluded from the entry-point scan.
func checkDomainAlias(pass *analysis.Pass) {
	for _, file := range pass.Files {
		if isTestFile(pass, file) {
			continue
		}
		for _, imp := range file.Imports {
			checkDomainImport(pass, imp)
		}
	}
}

func checkDomainImport(pass *analysis.Pass, imp *ast.ImportSpec) {
	path := strings.Trim(imp.Path.Value, `"`)
	if !strings.Contains(path, "/internal/domain/") {
		return
	}
	if imp.Name == nil || imp.Name.Name != "domain" {
		pass.Reportf(imp.Pos(), "command package: import the domain package with the \"domain\" alias")
	}
}
