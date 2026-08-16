// Package domain sits at myinternal/domain, where the tier marker appears
// WITHOUT its leading slash. It is not the domain tier and no alias is required
// of an import of it.
package domain

// Token is not a domain type.
type Token struct{ Value string }
