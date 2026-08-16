// Package config is a parent command holding a second verb it should not: the
// const-first check must anchor on THIS file, the one declaring the entry point
// whose verb is the package name, and not on whichever file sorts first among
// those declaring an entry point. apply.go sorts first and carries no const
// block, so anchoring by filename order files a metadata diagnostic against a
// file that carries no metadata and cannot satisfy the rule without being
// renamed.
package config

const name = "config"

// Command is the resident entry point: its verb is the package name.
func Command() string { return name }
