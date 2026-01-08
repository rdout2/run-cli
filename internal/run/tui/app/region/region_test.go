package region

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestRegionModal(t *testing.T) {
	app := tview.NewApplication()
	onSelect := func(region string) {}
	closeModal := func() {}

	modal := RegionModal(app, onSelect, closeModal)

	assert.NotNil(t, modal)
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok, "Expected RegionModal to return a Grid")
}
