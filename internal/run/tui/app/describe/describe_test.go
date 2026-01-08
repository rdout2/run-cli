package describe

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDescribeModal(t *testing.T) {
	app := tview.NewApplication()
	resource := map[string]string{"key": "value"}
	closeFunc := func() {
	}

	modal := DescribeModal(app, resource, "My Resource", closeFunc)

	assert.NotNil(t, modal)
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok, "Expected DescribeModal to return a Grid")
	
	// We can't easily trigger the InputCapture without simulating tcell events via Application,
	// which requires running the app. But we verified the structure creation.
}
