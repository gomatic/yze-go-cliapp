// Package picker is the near-miss of the keyword exemption, written against the
// MATCHER rather than the description. The exemption is keyed on token.Lookup,
// so a verb that merely RESEMBLES a keyword is not one: "selector" is not
// "select", the exemption does not apply, and the mismatch is reported exactly
// as any other would be.
package picker

const name = "picker"

// SelectorCommand names a verb that merely resembles a keyword.
func SelectorCommand() string { // want `the verb is the package name — move SelectorCommand into its own nested package "selector"`
	return name
}
