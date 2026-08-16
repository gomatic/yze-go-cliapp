// Package model is a type package nested beneath the nested domain group. A
// command naming one of its types imports it alongside its own counterpart,
// and it cannot also be called "domain" — the counterpart has that name.
package model

// Item is a nested domain type.
type Item struct{ ID string }
