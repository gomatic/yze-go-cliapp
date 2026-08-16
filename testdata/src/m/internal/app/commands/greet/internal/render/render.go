// Package render is a helper nested beneath the greet command package. Only
// the direct child of internal/app/commands is a command package; this deeper
// package is not one, so none of the command-package checks apply and it must
// produce no diagnostics.
package render

import wrong "m/internal/domain"

// Render deliberately violates every command-package rule (var-first file, no
// Command entry point, a domain import not aliased "domain") — each would be
// flagged if this nested package were misclassified as a command package.
var Render = func(cfg wrong.Argument) string { return string(cfg) }
