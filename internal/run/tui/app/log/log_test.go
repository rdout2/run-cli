package log

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestLogModal(t *testing.T) {
	app := tview.NewApplication()
	closeModal := func() {}

	modal := LogModal(app, "project", "filter", "My Logs", closeModal)

	assert.NotNil(t, modal)
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok, "Expected LogModal to return a Grid")
}
