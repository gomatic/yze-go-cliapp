// Package configure holds an entry point whose verb is a PREFIX of the package
// name. "config" is not "configure", so it is reported — widening the compare
// to HasPrefix silences a real finding and no fixture varied the two names'
// relationship until this one.
package configure

const name = "configure"

// ConfigCommand's verb is a prefix of the package name and not the package name.
func ConfigCommand() string { // want `the verb is the package name — move ConfigCommand into its own nested package "config"`
	return name
}
