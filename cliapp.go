// Package cliapp provides a go/analysis analyzer enforcing the app tier of the
// opinionated three-tier CLI layout, for every command package beneath
// internal/app/commands/ at any depth. A command package declares itself by
// exporting a command entry point; for each such package the command file (the
// one declaring the entry point whose verb is the package name) leads with a
// const block, every entry point's verb IS the package name — so a second verb
// belongs in its own nested package, as the reference layout demonstrates with
// config/{get,list,set} — and the package's own domain package is imported
// under the "domain" alias.
//
// The domain package is every import at or beneath an internal/domain tier,
// matched on whole path segments — so a sibling internal/domainhelpers is not
// the tier, and a module whose whole domain is a single package at
// internal/domain is. Two imports beneath it are exempt, each because the
// remedy would not compile: a blank import binds no identifier to spell
// "domain" with, and an import in a file where another domain-tier import
// already binds "domain" has no way to take that name — which is the ordinary
// shape of a command naming a type package nested under its own domain group.
//
// Test files are judged by none of this. They carry no command source, and a
// test naming several packages the command file does not may alias the domain
// import to disambiguate.
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
// A package that declares no entry point is out of scope entirely: it is a
// helper or grouping package, and depth never makes a package a command. The
// converse gap — a command that SHOULD exist but declares nothing — is not a
// single-package fact and is stickler/clilayout's to report.
//
// That declaration gate is also what excludes the two scaffolding passes the
// driver synthesizes for a package with tests: the external test package
// (clause "<pkg>_test") carries only _test.go files, which productionFiles
// drops, and the test-main package (import path "<pkg>.test") carries only the
// synthesized _testmain.go, which declares no entry point. Neither reaches a
// check, so neither needs a guard of its own — an earlier one was measured
// unable to change any verdict and removed. The remaining pass(es) hold the
// command source; the driver collapses their identical diagnostics, so the
// checks run once.
func run(pass *analysis.Pass) (any, error) {
	path := importPath(pass.Pkg.Path())
	if !isCommandTree(path) {
		return nil, nil
	}
	pkg := packageName(pass.Pkg.Name())
	commands := collectCommands(pass, pkg)
	if len(commands.funcs) == 0 {
		return nil, nil
	}
	reportNonConstFirst(pass, commands.file)
	checkVerbs(pass, pkg, commands.funcs)
	checkDomainAlias(pass)
	return nil, nil
}

// reportNonConstFirst reports when the command file's first non-import
// declaration is not a const block. The command file is the canonical metadata
// file, so the check targets it rather than an arbitrary first file of a
// multi-file package — see commandDecls for which file that is and why it is
// not decided by filename order.
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

// isTestFile reports whether file is a test file (_test.go). The name is taken
// from the token.File, which is the name the BUILD compiled the file under; the
// name go/token reports through a Position is whatever a //line directive in
// the file says it is, so a judged file could otherwise rename itself out of
// every check here.
func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go")
}

// productionFiles returns the package's non-test files, in file order.
func productionFiles(pass *analysis.Pass) []*ast.File {
	var found []*ast.File
	for _, file := range pass.Files {
		if !isTestFile(pass, file) {
			found = append(found, file)
		}
	}
	return found
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
// point in file order, and the command file — nil only when funcs is empty.
//
// The command file is the file declaring the RESIDENT entry point, the one
// whose verb is the package name. Filename order decides nothing here: the
// standard's canonical metadata file is the one carrying the package's own
// verb, and picking the alphabetically-first file that happens to declare an
// entry point files the const-first diagnostic against a file that carries no
// metadata, where the only edit that answers it is renaming files. When no entry
// point carries the package's verb the package is already reported for that
// (see checkVerbs) and the first declaring file is used, so the const-first
// rule still reaches a package whose verbs are all wrong.
type commandDecls struct {
	file  *ast.File
	funcs []*ast.FuncDecl
}

// collectCommands gathers the package's command declarations from its
// production files. Test files never carry the command source and are dropped
// by productionFiles: the in-package test-variant pass appends _test.go files
// to pass.Files, where an exported Test*Command test function (e.g.
// TestCommand) would otherwise be mistaken for an entry point — anchoring a
// spurious const-first diagnostic in the test file and marking a helper package
// as self-declaring. The result doubles as the self-declaration gate: no entry
// points marks a helper or grouping package.
func collectCommands(pass *analysis.Pass, pkg packageName) commandDecls {
	var found commandDecls
	declaredIn := map[*ast.FuncDecl]*ast.File{}
	for _, file := range productionFiles(pass) {
		for _, fn := range fileCommandFuncs(file) {
			found.funcs = append(found.funcs, fn)
			declaredIn[fn] = file
		}
	}
	if len(found.funcs) == 0 {
		return found
	}
	found.file = declaredIn[found.funcs[0]]
	if resident := residentCommand(found.funcs, pkg); resident != nil {
		found.file = declaredIn[resident]
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
// top-level function named Command or a <Verb>Command constructor. All three
// conjuncts narrow: a method named Command belongs to a type rather than to the
// package, an unexported one is invisible to the parent command that mounts it,
// and a name that merely CONTAINS "Command" (CommandLineFlags) is a helper.
// Widening any of them marks a helper package self-declaring and hands it a
// command package's obligations, which is the misclassification collectCommands
// exists to avoid. Recognition is otherwise deliberately broad — a package that
// exports PlanCommand has set out to be a command package and must be judged,
// not skipped — but its verb and its uniqueness are then checked; see
// checkVerbs.
func commandFunc(decl ast.Decl) (*ast.FuncDecl, bool) {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}
	return fn, fn.Recv == nil && ast.IsExported(fn.Name.Name) && strings.HasSuffix(fn.Name.Name, commandSuffix)
}
