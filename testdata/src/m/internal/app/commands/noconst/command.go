package noconst

import domain "m/internal/domain"

var cfg domain.Argument // want `const block`

func Command() domain.Argument { return cfg }
