package examples_test

// examples_test.go is the only Go file here and is an external test file, so the
// driver delivers a base-package pass with no syntax files. The import path
// still contains the command segment, so the analyzer must skip the empty pass
// instead of scanning it. Regression for the command-path examples-only crash.
func ExampleCommand() {}
