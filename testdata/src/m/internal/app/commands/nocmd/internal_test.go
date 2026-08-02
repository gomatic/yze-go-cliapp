package nocmd

import "testing"

// TestCommand produces the in-package test variant (command.go plus a _test.go
// file) and matches the exported *Command entry-point shape. The analyzer must
// skip _test.go files when collecting entry points, so this test function is
// never mistaken for one — mistaking it would mark this helper package as
// self-declaring and impose command obligations it does not carry.
func TestCommand(t *testing.T) { _ = t }
