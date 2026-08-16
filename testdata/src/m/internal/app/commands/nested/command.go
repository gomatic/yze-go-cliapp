// Package nested imports its own counterpart correctly AND a type package
// nested beneath it. "domain" is a file-scoped name and the counterpart has it;
// the type package is not the command's domain package and is not asked for a
// name it could not take.
package nested

import (
	domain "m/internal/domain/nested"
	"m/internal/domain/nested/model"
)

const name = "nested"

// Command is the entry point.
func Command() domain.Config { _ = model.Item{}; return domain.Config{Name: name} }
