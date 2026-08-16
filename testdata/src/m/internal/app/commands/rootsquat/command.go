// Package rootsquat is the other arrangement of the same two imports the vocab
// package gets right: the shared vocabulary at the tier root holds "domain"
// (its package clause is literally `package domain`, so an unaliased import
// binds it) and the command's own counterpart is left under another name.
// Reporting it with "alias it domain" would prescribe a redeclaration — the
// reversal that made the pair unsatisfiable in both directions — so the
// diagnostic names the holder, and renaming the vocabulary import is an edit
// that compiles.
package rootsquat

import (
	"m/internal/domain"

	wrong "m/internal/domain/rootsquat" // want `"m/internal/domain" already binds "domain"`
)

const name = "rootsquat"

// Command is the entry point.
func Command() string { return string(domain.Argument(wrong.Config{Name: name}.Name)) }
