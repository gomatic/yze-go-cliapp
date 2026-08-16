// Package rootbind imports the tier root under NO alias. The root's package
// clause is `package domain`, so the import already binds the name and every
// call site in the file reads domain.X exactly as the standard asks. Reading
// the alias rather than the name the import BINDS reports this file for an
// edit whose whole content is a redundant alias.
package rootbind

import "m/internal/domain"

const name = "rootbind"

// Command is the entry point.
func Command() domain.Argument { return name }
