package spinner

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	frames = []rune{
		'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏',
	}
	suffix = "Loading..."
)

// New returns a simple "⠋ Loading..." text, right-aligned, with animation.
func New(app *tview.Application) tview.Primitive {
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight)

	go func() {
		i := 0
		for {
			app.QueueUpdateDraw(func() {
				textView.SetText(fmt.Sprintf("%s %s", string(frames[i]), suffix))
			})
			i = (i + 1) % len(frames)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return textView
}
