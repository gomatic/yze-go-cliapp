// Package create is a CONFORMANT nested leaf command at depth two: const block
// first, a single Command() entry point, and the domain imported under the
// "domain" alias. Nothing is reported — and before the depth fix, nothing was
// even looked at.
package create

import domain "m/internal/domain/greet"

const name = "create"

// Command returns the CLI command definition.
func Command() string { return name + domain.Config{}.Name }
