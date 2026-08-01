package pkgstd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	pkgstd "github.com/gomatic/yze-go-pkgstd"
)

func TestCommandPackageStandards(t *testing.T) {
	analysistest.Run(
		t, analysistest.TestData(), pkgstd.Analyzer,
		"m/internal/domain/greet",
		"m/internal/app/commands/greet",
		"m/internal/app/commands/greet/internal/render",
		"m/internal/app/commands/multifile",
		"m/internal/app/commands/multiverb",
		"m/internal/app/commands/badalias",
		"m/internal/app/commands/noconst",
		"m/internal/app/commands/nocmd",
		"m/internal/app/commands/examples",
		"m/internal/app/commands/tested",
	)
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, pkgstd.Registration.Validate())
	assert.Equal(t, "yze/pkgstd", pkgstd.Registration.RuleID())
	assert.Same(t, pkgstd.Analyzer, pkgstd.Registration.Analyzer)
}

// TestCommandFileSkipsTestFilesWhenFindingTheEntryPoint names commandFile's
// claim. The in-package test-variant pass appends _test.go files to pass.Files,
// and an exported Test*Command function looks exactly like a command entry
// point to a naive scan. Mistaking one has two effects at once, and the second
// is the dangerous one: the const-first diagnostic is anchored in the TEST file
// (where it makes no sense), and the missing-entry-point diagnostic the real
// package deserved is masked — so a command package with no entry point at all
// passes the gate.
func TestCommandFileSkipsTestFilesWhenFindingTheEntryPoint(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), pkgstd.Analyzer, "m/...")
	require.NotEmpty(t, results)

	for _, r := range results {
		for _, d := range r.Diagnostics {
			position := r.Pass.Fset.Position(d.Pos)
			assert.NotContains(t, position.Filename, "_test.go",
				"a command-package diagnostic must never be anchored in a test file: %s", position)
		}
	}
}
