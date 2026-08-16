// Package blank mounts a domain package for its registration side effects. A
// blank import binds no identifier, so the "domain" alias has nothing to name
// and aliasing it leaves an import nothing uses — the exemption APPLIES and
// nothing is reported.
package blank

import _ "m/internal/domain/blank"

const name = "blank"

// Command is the entry point.
func Command() string { return name }
