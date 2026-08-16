package naming

// httptest.go ends in "test.go" and not in "_test.go", so the build compiles it
// as ordinary source. Matching the suffix without the underscore drops it.

import roundtrip "m/internal/domain/naming" // want `import the domain package with the "domain" alias`

var probe = roundtrip.Config{Name: "httptest"}
