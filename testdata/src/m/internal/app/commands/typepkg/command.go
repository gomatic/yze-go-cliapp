// Package typepkg names a TYPE PACKAGE nested beneath the domain tier and no
// counterpart of its own. A type package is not the command's domain package,
// so it is not asked for the name — the same scope limitation strayer states,
// one level deeper. It is the near-miss of vocab's silence: there the name is
// held by the counterpart, here there is no counterpart in the file at all, and
// both are silent for the same stated reason rather than by accident.
package typepkg

import item "m/internal/domain/nested/model"

const name = "typepkg"

// Command is the entry point.
func Command() item.Item { return item.Item{ID: name} }
