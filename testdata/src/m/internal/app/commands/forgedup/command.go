// Package forgedup tries to buy silence by importing the SAME path twice, once
// correctly aliased and once not, leaving every greetdomain. call site
// untouched. The second name is redundant rather than impossible — the file
// already has "domain" for that very package — so the exemption does not apply
// and the misaliased spelling is still reported.
package forgedup

import (
	domain "m/internal/domain/greet"
	greetdomain "m/internal/domain/greet" // want `import the domain package with the "domain" alias`
)

const name = "forgedup"

// Command is the entry point.
func Command() string { return name + domain.Config{}.Name + greetdomain.Config{}.Name }
