// Package render is a helper nested beneath the greet command package. Only a
// package declaring an entry point is a command package; this one declares
// none, so none of the command-package checks apply and it must produce no
// diagnostics.
//
// Its domain import is greet's, not its own: a nested helper's counterpart
// would be m/internal/domain/greet/internal/render, which Go's internal-
// visibility rule makes unimportable from here — so the alias rule can never
// reach a nested helper whatever the declaration gate does. The gate's own
// discrimination is carried by helper/, method/, midname/ and nocmd/, which sit
// directly beneath the tree and can import theirs.
package render

import greetdomain "m/internal/domain/greet"

// Render deliberately violates the command-package shape (a var-first file, no
// Command entry point) — each would be flagged if this nested package were
// misclassified as a command package that declares one.
var Render = func(cfg greetdomain.Config) string { return cfg.Name }
