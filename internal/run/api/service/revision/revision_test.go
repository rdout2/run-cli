package revision

import (
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMapRevision(t *testing.T) {
	now := time.Now()
	resp := &runpb.Revision{
		Name:       "projects/p/locations/l/services/s/revisions/my-rev",
		CreateTime: timestamppb.New(now),
		UpdateTime: timestamppb.New(now),
		Containers: []*runpb.Container{
			{
				Name:  "c1",
				Image: "img:latest",
				Resources: &runpb.ResourceRequirements{
					CpuIdle:         true,
					StartupCpuBoost: true,
				},
				Ports: []*runpb.ContainerPort{
					{
						ContainerPort: 8080,
					},
				},
			},
		},
		ExecutionEnvironment:          runpb.ExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN2,
		MaxInstanceRequestConcurrency: 80,
		Timeout:                       durationpb.New(time.Second * 30),
		NodeSelector: &runpb.NodeSelector{
			Accelerator: "nvidia-tesla-t4",
		},
	}

	result := mapRevision(resp, "my-service")

	assert.Equal(t, "my-rev", result.Name)
	assert.Equal(t, "my-service", result.Service)
	assert.Equal(t, now.Unix(), result.CreateTime.Unix())
	
	// Containers
	assert.Len(t, result.Containers, 1)
	assert.Equal(t, "c1", result.Containers[0].Name)
	assert.True(t, result.Containers[0].Resources.CPUIdle)
	
	// Env
	assert.Equal(t, "EXECUTION_ENVIRONMENT_GEN2", result.ExecutionEnvironment)
	assert.Equal(t, int32(80), result.MaxInstanceRequestConcurrency)
	assert.Equal(t, 30*time.Second, result.Timeout)
	
	// Accelerator
	assert.Equal(t, "nvidia-tesla-t4", result.Accelerator)
	
	// Top level shortcuts
	assert.True(t, result.CpuIdle)
	assert.True(t, result.StartupCpuBoost)
}
