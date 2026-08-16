// Package flatdomain imports the domain tier ROOT, which a module laying its
// domain out as one package has. The tier itself is in scope, so a wrong alias
// on it is reported — the boundary a "/internal/domain/" substring silently
// exempted.
package flatdomain

import wrong "m/internal/domain" // want `import the domain package with the "domain" alias`

const name = "flatdomain"

// Command is the entry point.
func Command() wrong.Argument { return name }
