package job

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

// Client defines the interface for Cloud Run Job operations.
type Client interface {
	ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error)
	RunJob(ctx context.Context, name string) (*runpb.Execution, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListJobs lists jobs for a project and region.
func (c *GCPClient) ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewJobsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &runpb.ListJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var jobs []*runpb.Job
	it := client.ListJobs(ctx, req)
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
		jobs = append(jobs, resp)
	}

	return jobs, nil
}

// RunJob runs a job.
func (c *GCPClient) RunJob(ctx context.Context, name string) (*runpb.Execution, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewJobsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	op, err := client.RunJob(ctx, &runpb.RunJobRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
