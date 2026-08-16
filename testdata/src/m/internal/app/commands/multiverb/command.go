package multiverb

// PlanCommand is a SECOND verb in the same package, declared in a second file —
// exactly the multi-file multi-verb package shape that shipped undetected.
func PlanCommand() string { // want `the verb is the package name — move PlanCommand into its own nested package "plan"`
	return name + " plan"
}

// helperCommand is unexported and must not satisfy the entry point.
func helperCommand() string { return name }
