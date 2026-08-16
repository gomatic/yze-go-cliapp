package cliapp

// The verb a command package answers to. The package doc states it as the
// standard — "the verb is the package name" — and everything that used to be
// decided by filename order is decided here instead: which entry point stays,
// which file is the canonical metadata file, and which symbol a diagnostic
// tells the author to move.

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// commandSuffix is the suffix every command entry point's name carries; what
// precedes it is the verb, and an empty remainder names the package's own verb.
const commandSuffix = "Command"

// verbOf returns the verb an entry point declares: what precedes the Command
// suffix, or the package's own name when the entry point is spelled Command.
func verbOf(fn *ast.FuncDecl, pkg packageName) string {
	if verb := strings.TrimSuffix(fn.Name.Name, commandSuffix); verb != "" {
		return strings.ToLower(verb)
	}
	return string(pkg)
}

// isPackageVerb reports whether an entry point's verb is the package name. The
// comparison folds case because the verb is spelled in an exported identifier
// and the package name is not: PlanCommand in package plan declares the verb
// the package is named for, and reporting it would be firing where the rule's
// own reason holds.
func isPackageVerb(fn *ast.FuncDecl, pkg packageName) bool {
	return strings.EqualFold(verbOf(fn, pkg), string(pkg))
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
// whose verb this same rule forbids.
func checkVerbs(pass *analysis.Pass, pkg packageName, commands []*ast.FuncDecl) {
	resident := residentCommand(commands, pkg)
	for _, fn := range commands {
		switch {
		case fn == resident:
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
