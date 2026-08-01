// Package pkgstd provides a go/analysis analyzer enforcing the gomatic three-tier
// layout's per-package standards. For a command package
// (internal/app/commands/<cmd>): the command file (the first one defining a
// command entry point) leads with a const block, at least one entry point —
// Command() or an exported <Verb>Command() constructor for a package exposing
// several verbs (e.g. PlanCommand/ApplyCommand) — exists, and the domain package
// is imported under the "domain" alias. (Cross-package correspondence is the
// layout analyzer's job.)
package pkgstd

import (
	"go/ast"
	"go/token"
	"strings"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
)

// Analyzer reports per-package violations of the three-tier command-package layout.
var Analyzer = &analysis.Analyzer{
	Name: "pkgstd",
	Doc:  "reports command packages that violate the gomatic three-tier package standards",
	Run:  run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "pkgstd",
	Categories: []goyze.Category{"structure"},
	URL:        "https://docs.gomatic.dev/yze/pkgstd",
	Analyzer:   Analyzer,
}

// run checks command packages for the per-package layout standards.
//
// The checks concern a command package's source code, so run skips the two
// scaffolding passes the driver synthesizes for a package that has tests — the
// external test package (clause "<pkg>_test") and the test-main package (import
// path "<pkg>.test") — neither of which carries the command source, so each
// would falsely trip the missing-entry-point check. The remaining pass(es) hold
// the command source; the driver collapses their identical diagnostics, so the
// checks run once. Skipping the scaffolding also guarantees a non-empty
// pass.Files (the test-only directory whose sole files are external tests yields
// an empty "<pkg>_test" pass), so checkCommandFunc never indexes an empty slice.
func run(pass *analysis.Pass) (any, error) {
	if isScaffoldingPackage(pass) || !isCommandPackage(importPath(pass.Pkg.Path())) || len(pass.Files) == 0 {
		return nil, nil
	}
	checkConstFirst(pass)
	checkCommandFunc(pass)
	checkDomainAlias(pass)
	return nil, nil
}

// checkConstFirst reports when the command file's first non-import declaration
// is not a const block. The command file (the first one defining a command
// entry point) is the canonical metadata file, so the check targets it rather
// than an arbitrary first file of a multi-file package. When no command file
// exists, checkCommandFunc reports the missing entry point and this check is a
// no-op.
func checkConstFirst(pass *analysis.Pass) {
	file := commandFile(pass)
	if file == nil {
		return
	}
	reportNonConstFirst(pass, file)
}

// reportNonConstFirst reports when file's first non-import declaration is not a
// const block.
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

// commandFile returns the first non-test file defining a command entry point,
// or nil when the package has none. Test files never carry the command source
// and are skipped: the in-package test-variant pass appends _test.go files to
// pass.Files, where an exported Test*Command test function (e.g. TestCommand)
// would otherwise be mistaken for the entry point — anchoring a spurious
// const-first diagnostic in the test file and masking a missing entry point.
func commandFile(pass *analysis.Pass) *ast.File {
	for _, file := range pass.Files {
		if !isTestFile(pass, file) && declaresCommandFunc(file) {
			return file
		}
	}
	return nil
}

// isTestFile reports whether file is a test file (_test.go).
func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go")
}

// declaresCommandFunc reports whether file declares a command entry point.
func declaresCommandFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if isCommandFunc(decl) {
			return true
		}
	}
	return false
}

func isImportDecl(decl ast.Decl) bool {
	gen, ok := decl.(*ast.GenDecl)
	return ok && gen.Tok == token.IMPORT
}

func isConstDecl(decl ast.Decl) bool {
	gen, ok := decl.(*ast.GenDecl)
	return ok && gen.Tok == token.CONST
}

// checkCommandFunc reports when the package has no command entry point.
func checkCommandFunc(pass *analysis.Pass) {
	if commandFile(pass) == nil {
		pass.Reportf(pass.Files[0].Name.Pos(), "command package: missing a Command() (or <Verb>Command()) entry point")
	}
}

// isCommandFunc reports whether decl is a command entry point: an exported
// top-level function named Command, or a <Verb>Command constructor — a package
// exposing several verbs from one domain (e.g. PlanCommand/ApplyCommand)
// satisfies the standard without a bare Command().
func isCommandFunc(decl ast.Decl) bool {
	fn, ok := decl.(*ast.FuncDecl)
	return ok && fn.Recv == nil && ast.IsExported(fn.Name.Name) && strings.HasSuffix(fn.Name.Name, "Command")
}

// checkDomainAlias reports domain imports not aliased as "domain".
func checkDomainAlias(pass *analysis.Pass) {
	for _, file := range pass.Files {
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
