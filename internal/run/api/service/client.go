package service

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client defines the interface for Cloud Run Service operations.
type Client interface {
	ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error)
	GetService(ctx context.Context, name string) (*runpb.Service, error)
	UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error)
}

// Ensure GCPClient implements Client
var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListServices lists services for a project and region.
func (c *GCPClient) ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &runpb.ListServicesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var services []*runpb.Service
	it := client.ListServices(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}
			return nil, err
		}
		services = append(services, resp)
	}

	return services, nil
}

// GetService gets a single service.
func (c *GCPClient) GetService(ctx context.Context, name string) (*runpb.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.GetService(ctx, &runpb.GetServiceRequest{Name: name})
}

// UpdateService updates a service.
func (c *GCPClient) UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	op, err := client.UpdateService(ctx, &runpb.UpdateServiceRequest{Service: service})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
