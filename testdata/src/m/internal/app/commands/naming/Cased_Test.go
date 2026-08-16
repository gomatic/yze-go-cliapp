package naming

// Cased_Test.go differs from a test file's name by one capital letter, and the
// go tool's rule is case-sensitive, so the build compiles it as ordinary
// source. Case-folding the name before matching the suffix drops it — the
// ordinary instinct of anyone who has been bitten by a case-insensitive
// filesystem, which is what this fleet develops on.

import cased "m/internal/domain/greet" // want `import the domain package with the "domain" alias`

var casedProbe = cased.Config{Name: "cased"}
