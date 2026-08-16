// Package nestedbad imports a package nested beneath the domain tier and
// nothing else from it, so no import in this file binds "domain". The exemption
// in `nested` is keyed on the NAME being taken, not on the import being nested,
// so it DOES NOT APPLY here and the import is reported.
package nestedbad

import item "m/internal/domain/greet/model" // want `import the domain package with the "domain" alias`

const name = "nestedbad"

// Command is the entry point.
func Command() item.Item { return item.Item{ID: name} }
