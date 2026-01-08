package project

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type ProjectsClientWrapper interface {
	SearchProjects(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper
	Close() error
}

type ProjectIteratorWrapper interface {
	Next() (*resourcemanagerpb.Project, error)
}

// Variables for dependency injection
var findDefaultCredentials = google.FindDefaultCredentials
var createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
	c, err := resourcemanager.NewProjectsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPProjectsClientWrapper{client: c}, nil
}

// Real implementations
type GCPProjectsClientWrapper struct {
	client *resourcemanager.ProjectsClient
}

func (w *GCPProjectsClientWrapper) SearchProjects(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
	return &GCPProjectIteratorWrapper{it: w.client.SearchProjects(ctx, req, opts...)}
}

func (w *GCPProjectsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPProjectIteratorWrapper struct {
	it *resourcemanager.ProjectIterator
}

func (w *GCPProjectIteratorWrapper) Next() (*resourcemanagerpb.Project, error) {
	return w.it.Next()
}

// Client defines the interface for Cloud Resource Manager operations.
type Client interface {
	ListProjects(ctx context.Context) ([]model.Project, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListProjects lists projects for the current user.
func (c *GCPClient) ListProjects(ctx context.Context) ([]model.Project, error) {
	// Explicitly find default credentials
	creds, err := findDefaultCredentials(ctx, resourcemanager.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	client, err := createProjectsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	req := &resourcemanagerpb.SearchProjectsRequest{
		// Query: "", // Empty query lists all projects
	}

	var projects []model.Project
	it := client.SearchProjects(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
				return nil, fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
			}
			return nil, err
		}

		projects = append(projects, mapProject(resp))
	}

	return projects, nil
}

func mapProject(resp *resourcemanagerpb.Project) model.Project {
	// Parse Project Number from Name "projects/123456"
	parts := strings.Split(resp.Name, "/")
	numberStr := parts[len(parts)-1]
	number, _ := strconv.Atoi(numberStr)

	return model.Project{
		Name:   resp.ProjectId,
		Number: number,
	}
}
