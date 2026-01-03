package revision

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	model "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List returns a list of revisions for the given service.
func List(project, region, service string) ([]model.Revision, error) {
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	c, err := run.NewRevisionsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		// Ignore error on close
		_ = c.Close()
	}()

	req := &runpb.ListRevisionsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, service),
	}

	var revisions []model.Revision
	it := c.ListRevisions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		nameParts := strings.Split(resp.Name, "/")
		name := nameParts[len(nameParts)-1]

		var cpuIdle, startupCpuBoost bool
		if len(resp.Containers) > 0 && resp.Containers[0].Resources != nil {
			cpuIdle = resp.Containers[0].Resources.CpuIdle
			startupCpuBoost = resp.Containers[0].Resources.StartupCpuBoost
		}

		revisions = append(revisions, model.Revision{
			Name:                          name,
			CreateTime:                    resp.CreateTime.AsTime(),
			UpdateTime:                    resp.UpdateTime.AsTime(),
			Service:                       service,
			ExecutionEnvironment:          resp.ExecutionEnvironment.String(),
			MaxInstanceRequestConcurrency: resp.MaxInstanceRequestConcurrency,
			Timeout:                       resp.Timeout.AsDuration(),
			CpuIdle:                       cpuIdle,
			StartupCpuBoost:               startupCpuBoost,
		})
	}

	return revisions, nil
}
