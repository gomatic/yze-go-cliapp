// Package dotimport imports its domain package into the file scope. Unlike a
// blank import, a dot import DOES bind identifiers — every exported name the
// domain package has — so "import it under the domain alias" is an edit the
// author can make, qualifying the uses as it goes. Exempting it would make one
// character the cheapest way to turn this rule off.
package dotimport

import . "m/internal/domain/greet" // want `import the domain package with the "domain" alias`

const name = "dotimport"

// Command is the entry point.
func Command() Config { return Config{Name: name} }
