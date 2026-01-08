package workerpool

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	app := tview.NewApplication()
	tbl := List(app)
	assert.NotNil(t, tbl)
	assert.Equal(t, LIST_PAGE_TITLE, tbl.Title)
}

func TestGetSelectedWorkerPoolFull(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	
	// Manually populate internal state
	workers = []model_workerpool.WorkerPool{
		{
			Name: "projects/p/locations/r/workerPools/pool-1",
			DisplayName: "pool-1",
			Region: "us-central1",
			UpdateTime: time.Now(),
			Scaling: &model_scaling.Scaling{ManualInstanceCount: 2},
		},
	}
	
	// Set selection to row 1
	listTable.Table.Select(1, 0)
	
	w := GetSelectedWorkerPoolFull()
	assert.NotNil(t, w)
	assert.Equal(t, "pool-1", w.DisplayName)
	
	// Header
	listTable.Table.Select(0, 0)
	w = GetSelectedWorkerPoolFull()
	assert.Nil(t, w)
}

func TestShortcuts(t *testing.T) {
	_ = header.New(info.Info{})
	
	assert.NotPanics(t, func() {
		Shortcuts()
	})
	
	assert.Contains(t, header.ContextShortcutView.GetText(true), "Refresh")
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	
	testWorkers := []model_workerpool.WorkerPool{
		{
			DisplayName: "pool-1",
			Region: "us-central1",
			Scaling: &model_scaling.Scaling{ManualInstanceCount: 5},
		},
	}
	
	render(testWorkers)
	
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "pool-1", listTable.Table.GetCell(1, 0).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 3).Text, "Manual: 5")
}
