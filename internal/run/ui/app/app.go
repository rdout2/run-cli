package app

import (
	"github.com/JulienBreux/run-cli/internal/run/auth"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_project "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/job"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/region"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/service"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/worker"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/spinner"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	pages *tview.Pages

	previousPageID string
	currentPageID  string
	currentInfo    info.Info

	projectModal tview.Primitive
	regionModal  tview.Primitive
	
	footerPages *tview.Pages
)

const (
	FULLSCREEN   = true
	ENABLE_MOUSE = true

	ESCAPE_SHORTCUT = tcell.KeyEscape
)

// ran the application.
func Run() error {
	app = tview.NewApplication()
	app.SetInputCapture(shortcuts)

	// Modals.
	projectModal = project.ProjectModal(app, func(selectedProject model_project.Project) {
		currentInfo.Project = selectedProject.Name
		header.UpdateInfo(currentInfo)
		switchTo(previousPageID)
	}, func() {
		switchTo(service.LIST_PAGE_ID)
	})

	regionModal = region.RegionModal(app, func(selectedRegion string) {
		currentInfo.Region = selectedRegion
		header.UpdateInfo(currentInfo)
		switchTo(previousPageID)
	}, func() {
		switchTo(service.LIST_PAGE_ID)
	})

	// Auth.

	// Load data.
	currentInfo = info.Info{
		User:    "Guest",
		Project: "None",
		Version: "dev",
		Region:  "-",
	}

	// Try to load real info
	if realInfo, err := auth.GetInfo(); err == nil {
		currentInfo.User = realInfo.User
		currentInfo.Project = realInfo.Project
		currentInfo.Region = realInfo.Region
	}

	return app.SetRoot(layout(), FULLSCREEN).
		EnableMouse(ENABLE_MOUSE).
		Run()
}

// returns the application layout.
func layout() *tview.Flex {
	pages = tview.NewPages()
	// Lists.
	pages.AddPage(service.LIST_PAGE_ID, service.List().Table, true, true)
	pages.AddPage(job.LIST_PAGE_ID, job.List().Table, true, true)
	pages.AddPage(worker.LIST_PAGE_ID, worker.List().Table, true, true)

	// Modals.
	pages.AddPage(project.MODAL_PAGE_ID, projectModal, true, true)
	pages.AddPage(region.MODAL_PAGE_ID, regionModal, true, true)

	// Footer (Spinner)
	footerPages = tview.NewPages()
	footerPages.AddPage("empty", tview.NewBox(), true, true)
	footerPages.AddPage("loading", spinner.New(), true, false)

	// pages.AddPage("details", detailView, true, false)
	// pages.AddPage("logs", logView, true, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header.New(currentInfo), 7, 1, false).
		// AddItem(titleBar, 1, 1, false).
		// AddItem(searchField, 1, 1, false). // New Search Bar (Height 1)
		AddItem(pages, 0, 1, true).
		AddItem(footerPages, 1, 1, false) // Footer

	// Default page
	switchTo(service.LIST_PAGE_ID)

	return layout
}

// shortcuts captures all key events.
func shortcuts(event *tcell.EventKey) *tcell.EventKey {
	// Navigation.
	if event.Key() == service.LIST_PAGE_SHORTCUT {
		switchTo(service.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == job.LIST_PAGE_SHORTCUT {
		switchTo(job.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == worker.LIST_PAGE_SHORTCUT {
		switchTo(worker.LIST_PAGE_ID)
		return nil
	}

	// Modals.
	if event.Key() == project.MODAL_PAGE_SHORTCUT {
		switchTo(project.MODAL_PAGE_ID)
		return nil
	}
	if event.Key() == region.MODAL_PAGE_SHORTCUT {
		switchTo(region.MODAL_PAGE_ID)
		return nil
	}

	// Open URL for Service list
	if currentPageID == service.LIST_PAGE_ID {
		if result := service.HandleShortcuts(event); result == nil {
			return nil
		}
	}

	return event
}

func showLoading() {
	footerPages.SwitchToPage("loading")
}

func hideLoading() {
	footerPages.SwitchToPage("empty")
}

func switchTo(pageID string) {
	previousPageID = currentPageID
	currentPageID = pageID
	pages.SwitchToPage(pageID)

	switch pageID {
	case service.LIST_PAGE_ID:
		service.Shortcuts()
		showLoading()
		service.ListReload(app, currentInfo, hideLoading)
	case job.LIST_PAGE_ID:
		job.Shortcuts()
		showLoading()
		job.ListReload(app, currentInfo, hideLoading)
	case worker.LIST_PAGE_ID:
		worker.Shortcuts()
		showLoading()
		worker.ListReload(app, currentInfo, hideLoading)
	case project.MODAL_PAGE_ID:
		header.ContextShortcutView.Clear()
		app.SetFocus(projectModal)
	case region.MODAL_PAGE_ID:
		header.ContextShortcutView.Clear()
		app.SetFocus(regionModal)
	}
}
