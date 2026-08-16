// Package squatter binds "domain" to a package that is not a domain package at
// all, leaving its own counterpart under another name. The name is taken, so
// "alias it domain" would prescribe a redeclaration; the diagnostic names the
// holder instead and the author's edit — rename the squatter, then alias the
// counterpart — is two steps that both compile.
package squatter

import (
	domain "strings"

	wrong "m/internal/domain/squatter" // want `already binds "domain"`
)

const name = "squatter"

// Command is the entry point.
func Command() string { return domain.ToUpper(wrong.Config{Name: name}.Name) }
