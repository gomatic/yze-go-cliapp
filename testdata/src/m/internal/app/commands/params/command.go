// Package params declares an entry point that TAKES PARAMETERS, which the real
// fleet does — kilroy's sessions and tui commands both do — and which every
// other fixture here holds constant at zero. Recognition is by receiver, export
// and suffix and by nothing else, so this package is self-declaring and carries
// every obligation: its var-first file and its misaliased domain import are
// both reported. Narrowing recognition to a zero-parameter signature takes the
// package out of scope entirely and silences both.
package params

import wrong "m/internal/domain/params" // want `import the domain package with the "domain" alias`

var defaults = wrong.Config{Name: "params"} // want `the first declaration must be the const block`

// Command is the entry point, and it takes its collaborators as parameters.
func Command(cfg wrong.Config) string { return cfg.Name + defaults.Name }
