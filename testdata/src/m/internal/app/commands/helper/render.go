// Package helper sits in the command tree and exports a function that is NOT a
// command entry point: its name does not end in Command. It declares no entry
// point, so it is a helper and carries none of the command package's
// obligations — although it violates every one of them (var-first file, a
// domain import under a wrong alias). Nothing is reported. Widening recognition
// to "any exported function" marks it self-declaring and reports both.
package helper

import wrong "m/internal/domain/greet"

var config = wrong.Config{Name: "helper"}

// Render is an ordinary exported helper.
func Render() string { return config.Name }
