package nocmd_test

// example_test.go is a black-box test file. The driver delivers it as an
// external test package (clause nocmd_test) plus a synthesized test-main
// package (import path nocmd.test); neither carries the command source, so
// isScaffoldingPackage skips both before any check runs.
func ExampleCommand() {}
