// Package mismatch declares exactly ONE entry point and still violates the
// standard: the verb is the package name, and "plan" is not "mismatch". The
// reference layout spells this commands/plan.
package mismatch

const name = "mismatch"

// PlanCommand names a verb the package is not named for.
func PlanCommand() string { return name } // want `the verb is the package name — move PlanCommand into its own nested package "plan"`
