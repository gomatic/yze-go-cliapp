// Package nested imports its domain group correctly AND a type package nested
// beneath it. "domain" is a file-scoped name already bound by the group, so the
// exemption APPLIES to the nested import and nothing is reported — asking for a
// second "domain" in this file prescribes a redeclaration.
package nested

import (
	domain "m/internal/domain/greet"
	"m/internal/domain/greet/model"
)

const name = "nested"

// Command is the entry point.
func Command() domain.Config { _ = model.Item{}; return domain.Config{Name: name} }
