// Package domain holds the vocabulary every domain package's Run contract
// shares. It sits AT the tier root, which the reference layout has and which a
// "/internal/domain/" substring test can never match.
package domain

// Argument is one positional command-line argument.
type Argument = string
