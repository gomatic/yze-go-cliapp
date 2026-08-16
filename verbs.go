package cliapp

// The verb a command package answers to. The package doc states it as the
// standard — "the verb is the package name" — and everything that used to be
// decided by filename order is decided here instead: which entry point stays,
// which file is the canonical metadata file, and which symbol a diagnostic
// tells the author to move.

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// commandSuffix is the suffix every command entry point's name carries; what
// precedes it is the verb, and an empty remainder names the package's own verb.
const commandSuffix = "Command"

// commandVerb is the CLI word a command package answers to — the thing the
// standard says the package is named for.
type commandVerb string

// verbOf returns the verb an entry point declares: what precedes the Command
// suffix, or the package's own name when the entry point is spelled Command.
func verbOf(fn *ast.FuncDecl, pkg packageName) commandVerb {
	if verb := strings.TrimSuffix(fn.Name.Name, commandSuffix); verb != "" {
		return commandVerb(strings.ToLower(verb))
	}
	return commandVerb(pkg)
}

// isPackageVerb reports whether an entry point's verb is the package name. The
// comparison folds case because the verb is spelled in an exported identifier
// and the package name is not: PlanCommand in package plan declares the verb
// the package is named for, and reporting it would be firing where the rule's
// own reason holds.
func isPackageVerb(fn *ast.FuncDecl, pkg packageName) bool {
	return strings.EqualFold(string(verbOf(fn, pkg)), string(pkg))
}

// isKeywordVerb reports whether a verb is a Go keyword, and therefore a name no
// package can have. The rule this file enforces is "the verb is the package
// name", and the go parser rejects `package import`, so an author whose CLI
// verb is import, select, go, type, range or map has already spelled the
// package as closely as Go permits (importer, selector) and has no edit left to
// make. Reporting it leaves the baseline as the remaining move, which is the
// population docs/r01.md exists to remove.
//
// Forging the exemption means naming a verb one of Go's keywords, which
// acquires the property it exempts — a verb no package can be named for —
// rather than a marker standing for it. TestKeywordVerbsAreTheOnesNoPackageCanHave
// pins both directions against token.Lookup.
func isKeywordVerb(verb commandVerb) bool {
	return token.Lookup(string(verb)).IsKeyword()
}

// residentCommand returns the entry point that legitimately stays in the
// package: the one spelled Command, else the first declaring the package's own
// verb, else nil. Preferring the bare spelling keeps the choice off filename
// order, which decides nothing about which verb a package is named for.
func residentCommand(commands []*ast.FuncDecl, pkg packageName) *ast.FuncDecl {
	var first *ast.FuncDecl
	for _, fn := range commands {
		switch {
		case !isPackageVerb(fn, pkg):
		case fn.Name.Name == commandSuffix:
			return fn
		case first == nil:
			first = fn
		}
	}
	return first
}

// checkVerbs reports every entry point that is not the package's resident one.
// One package is one verb, and the verb is the package name: the reference
// layout nests each verb in its own package (config/get, config/list,
// config/set), so a package exporting PlanCommand/ApplyCommand holds verbs that
// each belong in a nested package of their own, and a package named config
// whose sole entry point is PlanCommand is spelled config/plan.
//
// Naming the verb in the message is what makes the remedy takeable. Reporting
// every entry point after the first in FILE order names whichever symbol
// happened to sort late, which for a package holding Command and ApplyCommand
// is Command — and moving Command into a nested package produces config/command,
// whose verb this same rule forbids. A verb no package can be named for is
// exempt entirely; see isKeywordVerb.
func checkVerbs(pass *analysis.Pass, pkg packageName, commands []*ast.FuncDecl) {
	resident := residentCommand(commands, pkg)
	for _, fn := range commands {
		switch {
		case fn == resident:
		case isKeywordVerb(verbOf(fn, pkg)):
		case isPackageVerb(fn, pkg):
			pass.Reportf(
				fn.Pos(),
				"command package: one entry point per package — %s repeats the verb %q that %s already declares",
				fn.Name.Name, verbOf(fn, pkg), resident.Name.Name,
			)
		default:
			pass.Reportf(
				fn.Pos(),
				"command package: the verb is the package name — move %s into its own nested package %q",
				fn.Name.Name, verbOf(fn, pkg),
			)
		}
	}
}
