package scale

import (
	"testing"

	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestModal(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	
	service := &model_service.Service{
		Name: "s1",
		Scaling: &model_scaling.Scaling{
			ScalingMode: "AUTOMATIC",
			MinInstances: 1,
			MaxInstances: 5,
		},
	}
	
	onCompletion := func() {}
	
	modal := Modal(app, service, pages, onCompletion)
	
	assert.NotNil(t, modal)
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok, "Expected Modal to return a Grid")
}

func TestValidateScaleParams(t *testing.T) {
	tests := []struct {
		name      string
		mode      string
		manual    string
		min       string
		max       string
		wantErr   bool
		wantManual int64
		wantMin    int64
		wantMax    int64
	}{
		{"Valid Manual", "Manual", "5", "", "", false, 5, 0, 0},
		{"Invalid Manual", "Manual", "abc", "", "", true, 0, 0, 0},
		{"Valid Auto", "Automatic", "", "1", "5", false, 0, 1, 5},
		{"Valid Auto No Max", "Automatic", "", "1", "", false, 0, 1, 0},
		{"Invalid Auto Min", "Automatic", "", "abc", "5", true, 0, 0, 0},
		{"Invalid Auto Max", "Automatic", "", "1", "abc", true, 0, 0, 0},
		{"Invalid Auto Range", "Automatic", "", "5", "1", true, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, manual, err := validateScaleParams(tt.mode, tt.manual, tt.min, tt.max)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantManual, manual)
				assert.Equal(t, tt.wantMin, min)
				assert.Equal(t, tt.wantMax, max)
			}
		})
	}
}
