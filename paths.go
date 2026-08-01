package pkgstd

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

// isCommandPackage reports whether a package path is a command package: the
// direct child of internal/app/commands, i.e. exactly one path segment follows
// the marker. Deeper descendants (e.g. a helper nested beneath a command, such
// as .../commands/greet/internal/render) are not command packages and carry
// none of the per-package obligations.
func isCommandPackage(pkgPath importPath) bool {
	_, cmd, found := strings.Cut(string(pkgPath), "/internal/app/commands/")
	return found && cmd != "" && !strings.Contains(cmd, "/")
}

// isScaffoldingPackage reports whether pass is a driver-synthesized test
// package rather than a real package: an external test package (clause
// "<pkg>_test") or the test-main package (import path "<pkg>.test").
func isScaffoldingPackage(pass *analysis.Pass) bool {
	return strings.HasSuffix(pass.Pkg.Name(), "_test") || strings.HasSuffix(pass.Pkg.Path(), ".test")
}
