package cliapp

import "strings"

// Deciding which packages this analyzer judges, and which imports the alias
// rule judges. Everything here answers "is this in scope", never "is it
// correct" — a mistake in scope silently exempts a real command from every
// check that follows, and a green gate then asserts a standard nothing
// measured. Every predicate below therefore matches on whole path SEGMENTS: a
// substring match makes internal/domainhelpers a domain package and
// myinternal/app/commands a command tree, and neither mistake reports
// anything, so neither is visible in a diff of findings.

// commandsMarker is the path segment run below which a package may be a
// command package.
const commandsMarker = "/internal/app/commands/"

// importPath is the import path of an analyzed package or of one of its imports.
type importPath string

// packageName is the name a package's files declare in their package clause.
// It is the verb a command package answers to, which is why the entry point's
// own verb is compared against it (see isPackageVerb).
type packageName string

// isCommandTree reports whether a package path lies beneath internal/app/commands
// — at ANY depth, so a nested command like .../commands/tenant/create is in
// scope. Depth is not the discriminator between a command and a helper: whether
// the package declares a command entry point is (see run), so a helper nested
// beneath a command (.../commands/greet/internal/render) passes this predicate
// and is exempted by the declaration gate instead.
func isCommandTree(pkgPath importPath) bool {
	_, cmd, found := strings.Cut(string(pkgPath), commandsMarker)
	return found && cmd != ""
}
