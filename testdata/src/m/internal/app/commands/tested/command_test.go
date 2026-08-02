package tested

import (
	"testing"

	greetdomain "m/internal/domain/greet"
)

// TestCommand matches the exported *Command entry-point shape. In the
// in-package test-variant pass this file rides along in pass.Files; the
// analyzer must skip _test.go files when locating the command file so a test
// function is never mistaken for the entry point. The package is compliant and
// must produce no diagnostics.
func TestCommand(t *testing.T) {
	_ = Command()
	_ = t
}

// greetdomain is a disambiguating TEST alias for the domain import — legal in
// test files, which carry none of the command package's shape obligations.
var _ = greetdomain.Config{}
