package spinner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

// Spinner is a text view that displays a loading spinner.
type Spinner struct {
	*tview.TextView
	app     *tview.Application
	cancel  context.CancelFunc
	message string
	context string
	mu      sync.Mutex
}

// New returns a new spinner component.
func New(app *tview.Application) *Spinner {
	s := &Spinner{
		TextView: tview.NewTextView(),
		app:      app,
	}
	s.SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true).
		SetWrap(false)
	return s
}

// SetContext sets the contextual information displayed on the second line.
func (s *Spinner) SetContext(context string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.context = context
}

// Start starts the spinner animation with the given message.
func (s *Spinner) Start(message string) {
	s.Stop("") // Stop any existing animation

	s.mu.Lock()
	s.message = message
	s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	go func() {
		i := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.mu.Lock()
				msg := s.message
				ctxInfo := s.context
				s.mu.Unlock()

				s.app.QueueUpdateDraw(func() {
					text := fmt.Sprintf("%s %s", frames[i], msg)
					if ctxInfo != "" {
						text += fmt.Sprintf("\n[gray]%s", ctxInfo)
					}
					s.SetText(text)
				})
				i = (i + 1) % len(frames)
			}
		}
	}()
}

// Stop stops the spinner animation and sets the final message.
func (s *Spinner) Stop(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if message != "" {
		s.app.QueueUpdateDraw(func() {
			s.SetText(message)
		})
	}
}
