// Package domainhelpers is a sibling of the domain tier whose path merely
// BEGINS with the tier's. It is not the domain tier, and matching path
// segments rather than a substring is what tells them apart.
package domainhelpers

// Helper is not a domain type.
type Helper struct{ Name string }
