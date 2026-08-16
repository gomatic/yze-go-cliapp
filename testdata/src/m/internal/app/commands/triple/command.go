// Package triple names its counterpart THREE times: once correctly, twice not.
// EVERY other spec is reported, not merely the first — reporting one leaves the
// author an edit that does not end the finding, and a fixture with two specs
// cannot tell "all others" from "the first other".
package triple

import (
	alpha "m/internal/domain/triple" // want `already imported as "domain"`
	beta "m/internal/domain/triple"  // want `already imported as "domain"`
	domain "m/internal/domain/triple"
)

const name = "triple"

// Command is the entry point.
func Command() domain.Config {
	return domain.Config{Name: name + alpha.Config{}.Name + beta.Config{}.Name}
}
