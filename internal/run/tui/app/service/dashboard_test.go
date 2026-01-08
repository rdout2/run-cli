package service

import (
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_networking "github.com/JulienBreux/run-cli/internal/run/model/service/networking"
	model_security "github.com/JulienBreux/run-cli/internal/run/model/service/security"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	app := tview.NewApplication()
	d := Dashboard(app)
	assert.NotNil(t, d)
	// Header, Tabs, Pages
	assert.Equal(t, 3, d.GetItemCount())
}

func TestDashboardShortcuts(t *testing.T) {
	// Initialize header global view
	_ = header.New(info.Info{})
	
	assert.NotPanics(t, func() {
		DashboardShortcuts()
	})
	
	assert.Contains(t, header.ContextShortcutView.GetText(true), "Back")
}

func TestUpdateTabs(t *testing.T) {
	app := tview.NewApplication()
	_ = Dashboard(app)
	
	// Setup dummy service
	dashboardService = &model_service.Service{
		Name: "s1",
		Networking: &model_networking.Networking{
			Ingress: "INGRESS_TRAFFIC_ALL",
		},
		Security: &model_security.Security{
			InvokerIAMDisabled: true,
		},
	}
	
	assert.NotPanics(t, func() {
		updateNetworkingTab()
		assert.Contains(t, networkingDetail.GetText(true), "Allow all traffic")
		
		updateSecurityTab()
		assert.Contains(t, securityDetail.GetText(true), "Allow unauthenticated invocations")
	})
}
