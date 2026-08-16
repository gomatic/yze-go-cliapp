// Package vocabsplit imports its counterpart HERE and the shared vocabulary in
// its sibling file. Nothing is reported in either. The judgement is per file
// and the counterpart/vocabulary split is per package, so a fallback keyed on
// "this file names no counterpart" reports flags.go — and taking that remedy
// makes "domain" the counterpart in this file and the vocabulary in that one,
// inside one package.
package vocabsplit

import domain "m/internal/domain/vocabsplit"

const name = "vocabsplit"

// Command is the entry point.
func Command() domain.Config { return domain.Config{Name: name} }
