// Package forged binds the name "domain" to a package that is not the domain
// package. The exemption in `nested` exists because a second domain-tier import
// cannot hold a name the first already has; forging it here acquires the NAME
// and none of the property, so it does not apply and the real domain import is
// still reported. The author's remedy is to rename the squatter, which is an
// edit that exists — unlike the one a nested domain import would need.
package forged

import (
	domain "strings"

	wrong "m/internal/domain/greet" // want `import the domain package with the "domain" alias`
)

const name = "forged"

// Command is the entry point.
func Command() string { return domain.ToUpper(wrong.Config{Name: name}.Name) }
