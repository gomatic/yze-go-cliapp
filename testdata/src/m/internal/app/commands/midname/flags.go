// Package midname exports a function whose name CONTAINS "Command" without
// ending in it. It is not an entry point, so nothing is reported — although the
// file is var-first and its domain import carries a wrong alias. Matching the
// substring rather than the suffix marks it self-declaring and reports both.
package midname

import wrong "m/internal/domain/greet"

var config = wrong.Config{Name: "midname"}

// CommandLineFlags is an ordinary exported helper.
func CommandLineFlags() string { return config.Name }
