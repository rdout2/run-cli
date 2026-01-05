package command_test

import (
	"bytes"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/command"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	cmd := command.New(in, out, err)

	assert.NotNil(t, cmd)
	assert.Equal(t, "run", cmd.Use)
	assert.Equal(t, "Run is a CLI to play with Google Cloud Run interactively.", cmd.Short)
	
	// Check subcommands
	assert.True(t, cmd.HasSubCommands())
	
	// Check if version command exists
	found := false
	for _, c := range cmd.Commands() {
		if c.Use == "version" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected 'version' subcommand to be present")
}
