package cliapp_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	cliapp "github.com/gomatic/yze-go-cliapp"
)

func TestCommandPackageStandards(t *testing.T) {
	analysistest.Run(
		t, analysistest.TestData(), cliapp.Analyzer,
		"m/internal/domain/greet",
		"m/internal/domain",
		"m/internal/app/commands/greet",
		"m/internal/app/commands/greet/internal/render",
		"m/internal/app/commands/multifile",
		"m/internal/app/commands/multiverb",
		"m/internal/app/commands/badalias",
		"m/internal/app/commands/noconst",
		"m/internal/app/commands/nocmd",
		"m/internal/app/commands/examples",
		"m/internal/app/commands/tested",
		"m/internal/app/commands/tenant",
		"m/internal/app/commands/tenant/create",
		"m/internal/app/commands/tenant/badnested",
		"m/internal/app/commands/config",
		"m/internal/app/commands/mismatch",
		"m/internal/app/commands/dup",
		"m/internal/app/commands/blank",
		"m/internal/app/commands/nested",
		"m/internal/app/commands/typepkg",
		"m/internal/app/commands/lookalike",
		"m/internal/app/commands/unaliased",
		"m/internal/app/commands/helper",
		"m/internal/app/commands/midname",
		"m/internal/app/commands/method",
		"m/internal/app/commands/naming",
		"m/internal/app/commands/importer",
		"m/internal/app/commands/picker",
		"m/internal/app/commands/configure",
		"m/internal/app/commands/typefirst",
		"m/internal/app/commands/params",
		"m/internal/app/commands/generic",
		"m/internal/app/commands/rootsquat",
		"m/internal/app/commands/rawstring",
		"m/internal/app/commands/squatter",
		"m/internal/app/commands/vocab",
		"m/internal/app/commands/vocabonly",
		"m/internal/app/commands/vocabsplit",
		"m/internal/app/commands/domain",
		"m/internal/app/commands/triple",
		"m/internal/app/commands/twowrong",
		"m/internal/app/commands/strayer",
		"m/internal/app/commands/dotimport",
		"m/internal/app/commands/dupimport",
		"m/myinternal/app/commands/lookalike",
	)
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, cliapp.Registration.Validate())
	assert.Equal(t, "yze/cliapp", cliapp.Registration.RuleID())
	assert.Same(t, cliapp.Analyzer, cliapp.Registration.Analyzer)
}

// TestCollectCommandsNeverReadsTestFiles names collectCommands's claim that
// test files never carry the command source. The in-package test-variant pass
// appends _test.go files to pass.Files, and an exported Test*Command function
// looks exactly like a command entry point to a naive scan. Mistaking one has
// two effects at once: the const-first diagnostic is anchored in the TEST file
// (where it makes no sense), and a helper package is marked self-declaring —
// imposing command obligations it does not carry.
//
// A test file is one whose name ENDS in _test.go, which is the go tool's own
// rule and therefore the only one that answers "does the build run this". A
// substring test answers a different question and gets a different answer:
// helpers_test.golden.go is ordinary compiled source that contains "_test.go",
// and a diagnostic there is not a violation of anything this asserts.
func TestCollectCommandsNeverReadsTestFiles(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), cliapp.Analyzer, "m/...")
	require.NotEmpty(t, results)

	for _, r := range results {
		for _, d := range r.Diagnostics {
			position := r.Pass.Fset.Position(d.Pos)
			assert.False(t, strings.HasSuffix(position.Filename, "_test.go"),
				"a command-package diagnostic must never be anchored in a test file: %s", position)
		}
	}
}
