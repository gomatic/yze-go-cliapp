package naming

// helpers_test.golden.go CONTAINS "_test.go" without ending in it, and the
// build compiles it as ordinary source. Matching the substring anywhere in the
// name rather than at its end drops it.

import golden "m/internal/domain/naming" // want `import the domain package with the "domain" alias`

var goldenProbe = golden.Config{Name: "golden"}
