// Package unaliased imports its domain package under no alias at all, which is
// the commonest spelling of this violation and the one that leaves ImportSpec.Name
// nil.
package unaliased

import "m/internal/domain/greet" // want `import the domain package with the "domain" alias`

const name = "unaliased"

// Command is the entry point.
func Command() greet.Config { return greet.Config{Name: name} }
