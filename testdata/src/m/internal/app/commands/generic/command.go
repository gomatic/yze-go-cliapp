// Package generic declares a GENERIC entry point. Recognition is by receiver,
// export and suffix and by nothing else, so this package is self-declaring and
// carries every obligation — its var-first file is reported. Every other
// fixture holds the type-parameter dimension constant at "none", and narrowing
// recognition to reject a generic signature takes the package out of scope and
// silences it.
package generic

var prefix = "generic" // want `the first declaration must be the const block`

const name = "generic"

// Command is the entry point, and it is generic.
func Command[T ~string](v T) string { return prefix + name + string(v) }
