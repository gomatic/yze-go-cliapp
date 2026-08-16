package noconst

import domain "m/internal/domain/noconst"

var cfg domain.Config // want `const block`

func Command() domain.Config { return cfg }
