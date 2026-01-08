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
	flex, ok := l.(*tview.Flex)
	assert.True(t, ok, "Loader should return a *tview.Flex")
	assert.Equal(t, 4, flex.GetItemCount())
}
