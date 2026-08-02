// Package multiverb holds two verbs as two files in one package — the shape
// the standard forbids: one package is one verb, so each verb belongs in its
// own nested package (template.cli: config/{get,list,set}).
package multiverb

const name = "multiverb"

// ApplyCommand is the package's first entry point in file order; recognition
// marks the package self-declaring, and the first verb itself is not reported.
func ApplyCommand() string { return name + " apply" }
