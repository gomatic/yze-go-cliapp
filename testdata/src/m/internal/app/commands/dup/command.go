// Package dup declares two entry points that name the SAME verb — the package's
// own. Neither can be moved into a nested package named for its verb, because
// that verb is this package; one of them is simply redundant. The bare Command
// spelling is the resident, so which one is reported does not depend on
// filename order.
package dup

const name = "dup"

// Command is the resident entry point.
func Command() string { return name }

// DupCommand repeats the package's verb.
func DupCommand() string { return name } // want `one entry point per package — DupCommand repeats the verb "dup" that Command already declares`
