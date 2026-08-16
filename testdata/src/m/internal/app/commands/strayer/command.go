// Package strayer imports a domain group that is not its counterpart, under a
// wrong alias, and nothing is reported. That is a SCOPE LIMITATION stated in
// the doc comment rather than an exemption: "domain" names the package's own
// domain package, and greet is not it, so prescribing the name for greet would
// be firing where the rule's reason does not hold — and in a file that also
// named its counterpart, could not be obeyed at all. Whether a command may
// reach into another verb's domain package is cross-package correspondence,
// which is stickler/clilayout's.
package strayer

import greetdomain "m/internal/domain/greet"

const name = "strayer"

// Command is the entry point.
func Command() string { return name + greetdomain.Config{}.Name }
