package spinner

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// New returns a simple "Loading..." text, right-aligned.
func New() tview.Primitive {
	textView := tview.NewTextView().
		SetText("Loading ...").
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight)

	return textView
}
