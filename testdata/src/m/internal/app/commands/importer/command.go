// Package importer's CLI verb is "import", which no package can be named for —
// `package import` does not parse. The package is spelled as closely as Go
// permits, and there is no edit that makes the verb the package name, so the
// exemption APPLIES and nothing is reported. Reporting it would leave exactly
// one move, which is the baseline docs/r01.md exists to remove.
package importer

const name = "import"

// ImportCommand is the entry point for the import verb.
func ImportCommand() string { return name }
