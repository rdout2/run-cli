package job

import (
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMapJob(t *testing.T) {
	now := time.Now()
	resp := &runpb.Job{
		Name:    "projects/my-project/locations/us-central1/jobs/my-job",
		Creator: "user@example.com",
		LatestCreatedExecution: &runpb.ExecutionReference{
			Name:       "projects/my-project/locations/us-central1/executions/my-job-exec",
			CreateTime: timestamppb.New(now),
		},
		TerminalCondition: &runpb.Condition{
			State:              runpb.Condition_CONDITION_SUCCEEDED,
			Message:            "All good",
			LastTransitionTime: timestamppb.New(now),
		},
	}

	result := mapJob(resp, "us-central1")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Equal(t, "us-central1", result.Region)
	
	// Execution
	assert.NotNil(t, result.LatestCreatedExecution)
	assert.Equal(t, resp.LatestCreatedExecution.Name, result.LatestCreatedExecution.Name)
	assert.Equal(t, now.Unix(), result.LatestCreatedExecution.CreateTime.Unix())
	
	// Condition
	assert.NotNil(t, result.TerminalCondition)
	assert.Equal(t, "CONDITION_SUCCEEDED", result.TerminalCondition.State)
	assert.Equal(t, "All good", result.TerminalCondition.Message)
}
