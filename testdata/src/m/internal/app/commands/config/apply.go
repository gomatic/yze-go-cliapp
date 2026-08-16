package config

// ApplyCommand names the verb "apply", which belongs in config/apply. It is the
// symbol whose move is the fix, and it is the one the diagnostic names —
// reporting everything after the first in file order would name Command
// instead, whose only nested home would be config/command, a package whose verb
// the same rule forbids.
func ApplyCommand() string { // want `the verb is the package name — move ApplyCommand into its own nested package "apply"`
	return name + " apply"
}
