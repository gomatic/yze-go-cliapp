// Package dupimport buys the name without the property: it imports its own
// counterpart a SECOND time as "domain", uses it once, and leaves every call
// site reading dupdomain. The name is bound, so aliasing the second spec would
// redeclare it — the takeable edit is to drop the redundant import and use the
// one already called "domain", and that is what the message says. Without this
// the escape is two lines and the rule is off for the whole file.
package dupimport

import (
	domain "m/internal/domain/dupimport"
	dupdomain "m/internal/domain/dupimport" // want `already imported as "domain"`
)

const name = "dupimport"

var _ = domain.Config{}

// Command is the entry point.
func Command() dupdomain.Config { return dupdomain.Config{Name: name} }
