package nocmd_test

// example_test.go is a black-box test file. The driver delivers it as an
// external test package (clause nocmd_test) plus a synthesized test-main
// package (import path nocmd.test). Neither carries the command source: both
// hold only files productionFiles drops, so neither reaches a check and the
// declaration gate is what stops them. An earlier guard named them explicitly
// and was measured unable to change any verdict.
func ExampleCommand() {}
