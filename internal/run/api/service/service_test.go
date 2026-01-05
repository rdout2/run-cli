package service

import (
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMapService(t *testing.T) {
	// Setup input
	now := time.Now()
	resp := &runpb.Service{
		Name:         "projects/my-project/locations/us-central1/services/my-service",
		Uri:          "https://my-service.run.app",
		LastModifier: "user@example.com",
		UpdateTime:   timestamppb.New(now),
		Scaling: &runpb.ServiceScaling{
			MinInstanceCount: 1,
			MaxInstanceCount: 5,
			ScalingMode:      runpb.ServiceScaling_AUTOMATIC,
		},
		TrafficStatuses: []*runpb.TrafficTargetStatus{
			{
				Type:     runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST,
				Percent:  100,
				Revision: "my-service-v1",
			},
		},
		LatestReadyRevision:   "my-service-v1",
		LatestCreatedRevision: "my-service-v2",
		Ingress:               runpb.IngressTraffic_INGRESS_TRAFFIC_ALL,
	}

	// Call function
	result := mapService(resp, "my-project", "us-central1")

	// Assertions
	assert.Equal(t, "my-service", result.Name)
	assert.Equal(t, "https://my-service.run.app", result.URI)
	assert.Equal(t, "user@example.com", result.LastModifier)
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "us-central1", result.Region)
	
	// Scaling
	assert.Equal(t, "AUTOMATIC", result.Scaling.ScalingMode)
	assert.Equal(t, int32(1), result.Scaling.MinInstances)
	assert.Equal(t, int32(5), result.Scaling.MaxInstances)

	// Traffic
	assert.Len(t, result.TrafficStatuses, 1)
	assert.Equal(t, "my-service-v1", result.TrafficStatuses[0].Revision)
	
	// Revisions
	assert.Equal(t, "my-service-v1", result.LatestReadyRevision)
	assert.Equal(t, "my-service-v2", result.LatestCreatedRevision)

	// Networking
	assert.Equal(t, "INGRESS_TRAFFIC_ALL", result.Networking.Ingress)
}
