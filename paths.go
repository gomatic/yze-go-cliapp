package cliapp

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Deciding which packages this analyzer judges. Everything here answers
// "is this package in scope", never "is it correct" — a mistake in scope
// silently exempts a real command from every check that follows, and a green
// gate then asserts a standard nothing measured.

// importPath is the import path of an analyzed package.
type importPath string

// isCommandTree reports whether a package path lies beneath internal/app/commands
// — at ANY depth, so a nested command like .../commands/tenant/create is in
// scope. Depth is not the discriminator between a command and a helper: whether
// the package declares a command entry point is (see run), so a helper nested
// beneath a command (.../commands/greet/internal/render) passes this predicate
// and is exempted by the declaration gate instead.
func isCommandTree(pkgPath importPath) bool {
	_, cmd, found := strings.Cut(string(pkgPath), "/internal/app/commands/")
	return found && cmd != ""
}

// isScaffoldingPackage reports whether pass is a driver-synthesized test
// package rather than a real package: an external test package (clause
// "<pkg>_test") or the test-main package (import path "<pkg>.test").
func isScaffoldingPackage(pass *analysis.Pass) bool {
	return strings.HasSuffix(pass.Pkg.Name(), "_test") || strings.HasSuffix(pass.Pkg.Path(), ".test")
}
