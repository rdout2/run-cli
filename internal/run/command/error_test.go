package command_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/command"
	"github.com/stretchr/testify/assert"
)

func TestPrintError(t *testing.T) {
	w := &bytes.Buffer{}
	testErr := errors.New("something went wrong")

	err := command.PrintError(w, testErr)

	assert.NoError(t, err)
	assert.Equal(t, "something went wrong\n", w.String())
}
