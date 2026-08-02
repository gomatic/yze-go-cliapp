// Package badnested is a nested command package violating the const-first rule
// and the domain-alias rule — proving the checks REACH depth two, which the
// depth-1 predicate never did.
package badnested

import greetdomain "m/internal/domain/greet" // want `import the domain package with the "domain" alias`

var bad = greetdomain.Config{}.Name // want `the first declaration must be the const block`

// Command returns the CLI command definition.
func Command() string { return bad }
