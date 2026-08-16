// Package naming exists to vary the one dimension isTestFile reads: the file's
// NAME. Its three sibling files are ordinary source the build compiles and
// ships, each carrying a violation, and each is named so that a plausible
// loosening of the "_test.go" suffix test would drop it — case-folding the
// name, dropping the underscore, or matching the substring anywhere. A corpus
// that holds a matcher's input dimension constant cannot discriminate any
// widening of it.
package naming

import domain "m/internal/domain/naming"

const name = "naming"

// Command is the entry point.
func Command() domain.Config { return domain.Config{Name: name} }
