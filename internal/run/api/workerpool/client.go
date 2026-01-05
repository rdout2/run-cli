package workerpool

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

// Client defines the interface for Cloud Run WorkerPool operations.
type Client interface {
	ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error)
	GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error)
	UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListWorkerPools lists worker pools for a project and region.
func (c *GCPClient) ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &runpb.ListWorkerPoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var workerPools []*runpb.WorkerPool
	it := client.ListWorkerPools(ctx, req)
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
		workerPools = append(workerPools, resp)
	}

	return workerPools, nil
}

// GetWorkerPool gets a worker pool.
func (c *GCPClient) GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.GetWorkerPool(ctx, &runpb.GetWorkerPoolRequest{Name: name})
}

// UpdateWorkerPool updates a worker pool.
func (c *GCPClient) UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	op, err := client.UpdateWorkerPool(ctx, &runpb.UpdateWorkerPoolRequest{WorkerPool: workerPool})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
