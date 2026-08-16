// Package lookalike imports two packages whose paths resemble the domain tier
// and are not it: internal/domainhelpers, a sibling whose name merely begins
// with the tier's, and myinternal/domain, where the marker appears without its
// leading slash. Neither is a domain import, so no alias is required of either
// and nothing is reported.
package lookalike

import (
	helpers "m/internal/domainhelpers"
	token "m/myinternal/domain"
)

const name = "lookalike"

// Command is the entry point.
func Command() helpers.Helper { return helpers.Helper{Name: name + token.Token{Value: ""}.Value} }
