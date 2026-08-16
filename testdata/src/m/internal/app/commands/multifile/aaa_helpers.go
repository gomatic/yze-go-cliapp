package multifile

// aaa_helpers.go sorts before command.go, so it is pass.Files[0]. Its first
// declaration is deliberately a var, not a const: the const-first check must
// target the command file (command.go), not this arbitrarily-first file, so
// this file must draw no const-first diagnostic.
var helper = 1
