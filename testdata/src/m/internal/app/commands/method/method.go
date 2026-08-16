// Package method declares Command as a METHOD, which belongs to a type rather
// than to the package and cannot be mounted by a parent command. It is not an
// entry point, so nothing is reported — although the file is var-first and its
// domain import carries a wrong alias. Dropping the receiver conjunct marks the
// package self-declaring and reports both.
package method

import wrong "m/internal/domain/method"

type builder struct{}

var config = wrong.Config{Name: "method"}

// Command is a method, not a package entry point.
func (builder) Command() string { return config.Name }
