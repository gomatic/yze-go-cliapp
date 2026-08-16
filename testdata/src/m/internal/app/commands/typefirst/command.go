// Package typefirst's command file leads with a TYPE declaration. The rule is
// that the first declaration is the const block, so a leading type is reported
// exactly as a leading var is; widening the check to accept a type block
// silences a real finding, and every other fixture leads with a var.
package typefirst

// Options is the leading declaration, and it is not a const block.
type Options struct{ Name string } // want `the first declaration must be the const block`

const name = "typefirst"

// Command is the entry point.
func Command() Options { return Options{Name: name} }
