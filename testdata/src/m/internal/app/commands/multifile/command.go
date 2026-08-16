package multifile

import domain "m/internal/domain/multifile"

const (
	name  = "multifile"
	usage = "multi-file command"
)

// Command is the entry point; its file leads with the const block.
func Command() domain.Config {
	_ = name
	_ = usage
	return domain.Config{}
}
