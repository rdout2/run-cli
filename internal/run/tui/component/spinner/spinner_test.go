package spinner

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := tview.NewApplication()
	s := New(app)
	assert.NotNil(t, s)
	assert.NotNil(t, s.TextView)
	assert.Equal(t, app, s.app)
}

func TestStartStop(t *testing.T) {
	// Use SimulationScreen to allow app.Run() without terminal
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}

	app := tview.NewApplication()
	app.SetScreen(screen)

	// Run app in goroutine to process QueueUpdateDraw events
	go func() {
		if err := app.Run(); err != nil {
			// app.Run() might return error when stopped, which is expected
		}
	}()
	
	// Ensure app stops at end of test
	defer app.Stop()

	s := New(app)

	// Start
	s.Start("Loading...")
	
	s.mu.Lock()
	assert.NotNil(t, s.cancel, "Cancel function should be set after Start")
	s.mu.Unlock()

	// Wait a bit to let goroutine run and process updates
	time.Sleep(150 * time.Millisecond)

	// Stop
	s.Stop("Done")

	s.mu.Lock()
	assert.Nil(t, s.cancel, "Cancel function should be nil after Stop")
	s.mu.Unlock()
}