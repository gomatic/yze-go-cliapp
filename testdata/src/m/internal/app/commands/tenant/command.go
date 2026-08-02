// Package tenant is a mount-only parent command: it declares Command() to
// mount its nested verb packages and imports no domain. Nothing is reported —
// a parent that only mounts subcommands has no domain counterpart and no
// Action.
package tenant

const name = "tenant"

// Command mounts the nested verb subcommands.
func Command() string { return name }
