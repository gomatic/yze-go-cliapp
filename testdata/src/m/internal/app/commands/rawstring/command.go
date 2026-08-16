// Package rawstring writes its import path as a RAW string literal, which Go
// admits as readily as the interpreted form and which the build treats
// identically. Taking the quotes off by trimming the double quote leaves the
// backquotes on the path, and a path carrying them matches no predicate here —
// a domain import turned off by one keystroke. It is reported.
package rawstring

import wrong `m/internal/domain` // want `import the domain package with the "domain" alias`

const name = "rawstring"

// Command is the entry point.
func Command() wrong.Argument { return name }
