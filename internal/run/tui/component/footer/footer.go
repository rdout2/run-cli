package footer

import (
	"github.com/rivo/tview"
)

var (
	ContextShortcutView *tview.TextView
)

// New returns a TView footer.
// Footer composition:
// | Contextual Shortcuts |
func New() *tview.Flex {
	ContextShortcutView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	return tview.NewFlex().
		AddItem(ContextShortcutView, 0, 1, false)
}