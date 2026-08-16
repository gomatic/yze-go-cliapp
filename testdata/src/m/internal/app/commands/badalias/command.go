package badalias

import (
	greet "m/internal/domain/badalias" // want `domain.*alias`
)

const x = 1

func Command() greet.Config { return greet.Config{} }
