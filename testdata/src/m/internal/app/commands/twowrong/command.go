// Package twowrong names its counterpart twice and neither spec binds the
// standard's name. The finding sits on the FIRST, which is what makes the
// remedy end it: aliasing the first leaves the second a duplicate to drop,
// while reporting the last would prescribe an alias the first still holds.
package twowrong

import (
	alpha "m/internal/domain/twowrong" // want `import the domain package with the "domain" alias`
	beta "m/internal/domain/twowrong"
)

const name = "twowrong"

// Command is the entry point.
func Command() alpha.Config { return alpha.Config{Name: name + beta.Config{}.Name} }
