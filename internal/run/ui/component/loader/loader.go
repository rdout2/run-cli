package loader

import (
	"github.com/rivo/tview"
)

// New returns a new loader component.
func New() tview.Primitive {
	modal := tview.NewModal().
		SetText("Loading... Please wait").
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	return modal
}
