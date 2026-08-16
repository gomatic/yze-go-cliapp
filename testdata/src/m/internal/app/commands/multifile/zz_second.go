package multifile

// zz_second.go sorts LAST, which is the point: the alias scan runs over every
// production file, and a multi-file command package's later files are exactly
// where an unaliased domain import survives review. Narrowing the scan to
// pass.Files[0] leaves this finding unreported, and only an assertion — never a
// coverage percentage — says the rule stopped reaching the file it exists to
// reach.

import second "m/internal/domain/multifile" // want `import the domain package with the "domain" alias`

var helperConfig = second.Config{Name: "multifile"}
