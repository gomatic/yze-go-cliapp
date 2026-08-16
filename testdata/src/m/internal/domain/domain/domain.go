// Package domain is the domain counterpart of the "domain" command — a verb a
// registrar or DNS CLI really has. Its directory is domain, so its package
// clause is domain, so an UNALIASED import of it binds the standard's name
// with no alias to spell.
package domain

// Config is the command configuration.
type Config struct{ Name string }
