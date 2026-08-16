// Package lookalike sits at myinternal/app/commands, a path that CONTAINS the
// command-tree marker without carrying it as a leading segment. It is not a
// command package, so nothing is reported — although it violates every rule
// this analyzer has.
package lookalike

import wrong "m/internal/domain/greet"

var config = wrong.Config{Name: "lookalike"}

// Command would be the entry point if this were a command package.
func Command() string { return config.Name }
