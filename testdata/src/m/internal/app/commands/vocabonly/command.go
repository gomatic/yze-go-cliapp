// Package vocabonly is a mount-only parent command: it mounts nested verbs and
// has no domain counterpart of its own, while still naming the shared
// vocabulary package because its entry point takes domain arguments. Nothing is
// reported. Standing the tier root in for the counterpart a file does not
// import reports THIS, and the remedy — alias the vocabulary "domain" — makes
// the name mean the vocabulary in a command package whose sibling verbs mean
// their counterparts by it.
package vocabonly

import shared "m/internal/domain"

const name = "vocabonly"

// Command is the entry point.
func Command(args ...shared.Argument) string { return name + string(len(args)) }
