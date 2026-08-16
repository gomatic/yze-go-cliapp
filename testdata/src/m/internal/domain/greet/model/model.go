// Package model is a type package nested beneath the greet domain group. A
// command naming one of its types imports it alongside the group itself, and
// it cannot also be called "domain".
package model

// Item is a nested domain type.
type Item struct{ ID string }
