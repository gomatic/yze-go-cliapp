// Package forgeroot tries to buy silence by binding "domain" to the tier ROOT,
// which is the shared vocabulary package and nobody's domain counterpart. The
// alias is acquired and none of the property is, so the exemption does not
// apply and the file's own misaliased group is still reported. Excusing it here
// would make two lines the cheapest way to turn this rule off.
package forgeroot

import (
	domain "m/internal/domain"
	greetdomain "m/internal/domain/greet" // want `import the domain package with the "domain" alias`
)

const name = "forgeroot"

// Command is the entry point.
func Command() string { return name + greetdomain.Config{Name: string(domain.Argument(""))}.Name }
