package service

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type ServicesClientWrapper interface {
	ListServices(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper
	GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error)
	UpdateService(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error)
	Close() error
}

type ServiceIteratorWrapper interface {
	Next() (*runpb.Service, error)
}

type UpdateServiceOperationWrapper interface {
	Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error)
}

// Variables for dependency injection
var findDefaultCredentials = google.FindDefaultCredentials
var createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
	c, err := run.NewServicesClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPServicesClientWrapper{client: c}, nil
}

// Real implementations
type GCPServicesClientWrapper struct {
	client *run.ServicesClient
}

func (w *GCPServicesClientWrapper) ListServices(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper {
	return &GCPServiceIteratorWrapper{it: w.client.ListServices(ctx, req, opts...)}
}

func (w *GCPServicesClientWrapper) GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
	return w.client.GetService(ctx, req, opts...)
}

func (w *GCPServicesClientWrapper) UpdateService(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error) {
	op, err := w.client.UpdateService(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPUpdateServiceOperationWrapper{op: op}, nil
}

func (w *GCPServicesClientWrapper) Close() error {
	return w.client.Close()
}

type GCPServiceIteratorWrapper struct {
	it *run.ServiceIterator
}

func (w *GCPServiceIteratorWrapper) Next() (*runpb.Service, error) {
	return w.it.Next()
}

type GCPUpdateServiceOperationWrapper struct {
	op *run.UpdateServiceOperation
}

func (w *GCPUpdateServiceOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error) {
	return w.op.Wait(ctx, opts...)
}

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
	creds, err := findDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := createServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

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
	creds, err := findDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := createServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	return client.GetService(ctx, &runpb.GetServiceRequest{Name: name})
}

// UpdateService updates a service.
func (c *GCPClient) UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
	creds, err := findDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := createServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	op, err := client.UpdateService(ctx, &runpb.UpdateServiceRequest{Service: service})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
