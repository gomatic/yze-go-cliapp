// Package domain names its counterpart with no alias at all, and the
// counterpart's package clause is "domain" — so the import already binds the
// standard's name and every call site reads domain.Config. Nothing is reported.
// Reading the ALIAS rather than the name the import BINDS reports this file for
// an edit whose whole content is a redundant alias.
package domain

import "m/internal/domain/domain"

const name = "domain"

// Command is the entry point.
func Command() domain.Config { return domain.Config{Name: name} }
