package revision

import (
	"context"
	"fmt"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client defines the interface for Cloud Run Revision operations.
type Client interface {
	ListRevisions(ctx context.Context, project, region, service string) ([]*runpb.Revision, error)
}

var apiClient Client = &GCPClient{}

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListRevisions lists revisions for a service.
func (c *GCPClient) ListRevisions(ctx context.Context, project, region, service string) ([]*runpb.Revision, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := run.NewRevisionsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &runpb.ListRevisionsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, service),
	}

	var revisions []*runpb.Revision
	it := client.ListRevisions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, resp)
	}

	return revisions, nil
}
