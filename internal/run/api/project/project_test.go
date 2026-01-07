package project

import (
	"context"
	"errors"
	"testing"

	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/stretchr/testify/assert"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	ListProjectsFunc func(ctx context.Context) ([]model.Project, error)
}

func (m *MockClient) ListProjects(ctx context.Context) ([]model.Project, error) {
	return m.ListProjectsFunc(ctx)
}

func TestList(t *testing.T) {
	// Backup original client and restore after test
	origClient := apiClient
	defer func() { apiClient = origClient }()

	t.Run("success", func(t *testing.T) {
		expectedProjects := []model.Project{
			{Name: "project-1", Number: 123},
			{Name: "project-2", Number: 456},
		}

		apiClient = &MockClient{
			ListProjectsFunc: func(ctx context.Context) ([]model.Project, error) {
				return expectedProjects, nil
			},
		}

		projects, err := List()
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("something went wrong")

		apiClient = &MockClient{
			ListProjectsFunc: func(ctx context.Context) ([]model.Project, error) {
				return nil, expectedError
			},
		}

		projects, err := List()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, projects)
	})
}
