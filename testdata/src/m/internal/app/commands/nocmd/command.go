// Package nocmd sits in the command tree but declares no command entry point,
// so it does not self-declare and carries none of the per-package obligations —
// nothing is reported, including its unaliased-adjacent shape. Whether a
// command SHOULD exist here is cross-package correspondence, which is
// stickler/clilayout's to report.
package nocmd

import domain "m/internal/domain/greet"

const x = 1

var _ = domain.Config{}
