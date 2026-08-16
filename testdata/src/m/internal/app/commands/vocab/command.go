// Package vocab is the shape the standard itself produces and the suite could
// not accept: a command file naming its own domain counterpart AND the shared
// vocabulary package at the tier root, which yze/clidomain requires every
// domain Run signature to name. "domain" is file-scoped, the counterpart has
// it, and nothing else in the file can — so nothing else in the file is asked
// for it. Nothing is reported. Widening the rule back to every import at or
// beneath the tier reports the vocabulary import and prescribes a
// redeclaration, which is the only instruction an author cannot take.
package vocab

import (
	shared "m/internal/domain"
	domain "m/internal/domain/vocab"
)

const name = "vocab"

// Command is the entry point.
func Command(args ...shared.Argument) domain.Config {
	return domain.Config{Name: name + string(shared.Argument(len(args)))}
}
