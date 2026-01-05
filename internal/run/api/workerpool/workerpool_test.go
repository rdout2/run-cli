package workerpool

import (
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMapWorkerPool(t *testing.T) {
	now := time.Now()
	instanceCount := int32(2)
	resp := &runpb.WorkerPool{
		Name:         "projects/my-project/locations/us-central1/workerPools/my-pool",
		LastModifier: "user@example.com",
		UpdateTime:   timestamppb.New(now),
		Scaling: &runpb.WorkerPoolScaling{
			ManualInstanceCount: &instanceCount,
		},
		Labels: map[string]string{
			"env": "prod",
		},
	}

	result := mapWorkerPool(resp, "my-project", "us-central1")

	assert.Equal(t, "my-pool", result.DisplayName)
	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.LastModifier)
	assert.Equal(t, now.Unix(), result.UpdateTime.Unix())
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "us-central1", result.Region)
	
	// Scaling
	assert.NotNil(t, result.Scaling)
	assert.Equal(t, int32(2), result.Scaling.ManualInstanceCount)
	
	// Labels
	assert.Equal(t, "prod", result.Labels["env"])
}
