// Package multiverb holds two verbs as two files in one package — the shape
// the standard forbids: one package is one verb and the verb is the package
// name, so each of these belongs in a nested package named for it
// (template.cli: config/{get,list,set}). NEITHER verb is "multiverb", so
// neither is the resident and both are reported, each told the nested package
// its own verb names.
package multiverb

const name = "multiverb"

// ApplyCommand is the package's first entry point in file order. File order
// decides nothing: it is reported because its verb is not the package name.
func ApplyCommand() string { // want `the verb is the package name — move ApplyCommand into its own nested package "apply"`
	return name + " apply"
}
