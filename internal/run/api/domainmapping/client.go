package domainmapping

import (
	"context"
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

// Client defines the interface for the DomainMapping API client.
type Client interface {
	ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error)
}

// GCPClient is the Google Cloud Platform implementation of the Client interface.
type GCPClient struct{}

// ListDomainMappings lists domain mappings for a given project and region.
func (c *GCPClient) ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	runService, err := run.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create run service: %w", err)
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)

	var domainMappings []*run.DomainMapping
	pageToken := ""

	for {
		call := runService.Projects.Locations.Domainmappings.List(parent)
		if pageToken != "" {
			call.Continue(pageToken)
		}

		resp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to list domain mappings: %w", err)
		}

		domainMappings = append(domainMappings, resp.Items...)

		if resp.Metadata != nil {
			pageToken = resp.Metadata.Continue
		} else {
			pageToken = ""
		}
		
		if pageToken == "" {
			break
		}
	}

	return domainMappings, nil
}
