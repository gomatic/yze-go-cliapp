package multiverb

// PlanCommand is a SECOND verb in the same package, declared in a second file —
// exactly the multi-file multi-verb package shape that shipped undetected.
func PlanCommand() string { return name + " plan" } // want `one verb per package`

// helperCommand is unexported and must not satisfy the entry point.
func helperCommand() string { return name }
