package version_test

import (
	"bytes"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/command/version"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdVersion(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	cmd := version.NewCmdVersion(in, out, err)

	assert.NotNil(t, cmd)
	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Print the Run CLI version", cmd.Short)
	assert.Equal(t, "Print the Run CLI version", cmd.Long)

	// Check flags
	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, "One of '', 'yaml' or 'json'.", outputFlag.Usage)

	// Test execution
	cmd.SetOut(out)
	cmd.SetErr(err)
	
	// Execute the command
	execErr := cmd.Execute()
	assert.NoError(t, execErr)
	
	// Verify output contains version info
	// Note: exact content depends on pkg/version globals, but we expect at least "Version:"
	assert.Contains(t, out.String(), "Version:")
}
