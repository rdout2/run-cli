package project

import (
	"context"

	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
)

var apiClient Client = &GCPClient{}

// List returns a list of projects for the current user.
func List() ([]model.Project, error) {
	ctx := context.Background()
	return apiClient.ListProjects(ctx)
}