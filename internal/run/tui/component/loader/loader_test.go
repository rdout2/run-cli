package loader

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := tview.NewApplication()
	l := New(app)

	assert.NotNil(t, l)
	assert.NotNil(t, l.Flex)
	assert.NotNil(t, l.Spinner)
	assert.Equal(t, 4, l.GetItemCount())
}