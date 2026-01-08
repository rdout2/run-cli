package spinner

import (
	"testing"

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

// func TestStartStop(t *testing.T) {
// 	app := tview.NewApplication()
// 	s := New(app)

// 	// Start
// 	s.Start("Loading...")

// 	s.mu.Lock()
// 	assert.NotNil(t, s.cancel, "Cancel function should be set after Start")
// 	s.mu.Unlock()

// 	// Wait a bit to let goroutine run (though we can't verify drawing without UI loop)
// 	time.Sleep(150 * time.Millisecond)

// 	// Stop
// 	s.Stop("Done")

// 	s.mu.Lock()
// 	assert.Nil(t, s.cancel, "Cancel function should be nil after Stop")
// 	s.mu.Unlock()
// }
