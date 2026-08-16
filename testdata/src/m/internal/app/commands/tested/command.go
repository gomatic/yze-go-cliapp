package tested

import domain "m/internal/domain/tested"

const (
	name  = "tested"
	usage = "command package with an in-package test"
)

// Command is the entry point; its file leads with the const block.
func Command() domain.Config {
	_ = name
	_ = usage
	return domain.Config{}
}
