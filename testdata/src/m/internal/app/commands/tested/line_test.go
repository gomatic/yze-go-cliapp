//line aimed.go:1

// line_test.go is a REAL test file that names itself something else. A //line
// directive changes every position go/token reports for this file without
// changing anything the build does with it, so a file identity taken from a
// Position is a name the judged file wrote itself.
//
// TestPlanCommand carries the exported *Command shape deliberately. Read the
// compiled name and this file is a test file, skipped, and nothing is reported.
// Read the Position name and it is "aimed.go", ordinary source — TestPlanCommand
// becomes a second entry point whose verb is not the package name, and the
// diagnostic lands in a file that does not exist.
package tested

import "testing"

// TestPlanCommand keeps the package's real test honest about the shape.
func TestPlanCommand(t *testing.T) { _ = t }
