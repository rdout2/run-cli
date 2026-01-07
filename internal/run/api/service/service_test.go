package service

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
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

func TestMapService_Complex(t *testing.T) {
	instanceCount := int32(3)
	resp := &runpb.Service{
		Name: "projects/p/locations/r/services/s-complex",
		Scaling: &runpb.ServiceScaling{
			ScalingMode:         runpb.ServiceScaling_MANUAL,
			ManualInstanceCount: &instanceCount,
		},
		Template: &runpb.RevisionTemplate{
			VpcAccess: &runpb.VpcAccess{
				Connector: "connector1",
				Egress:    runpb.VpcAccess_ALL_TRAFFIC,
			},
			ServiceAccount: "sa@example.com",
			EncryptionKey:  "key1",
		},
		BinaryAuthorization: &runpb.BinaryAuthorization{
			BinauthzMethod: &runpb.BinaryAuthorization_UseDefault{
				UseDefault: true,
			},
			BreakglassJustification: "emergency",
		},
	}

	result := mapService(resp, "p", "r")

	// Scaling
	assert.Equal(t, "MANUAL", result.Scaling.ScalingMode)
	assert.Equal(t, int32(3), result.Scaling.ManualInstanceCount)

	// Networking / VPC
	assert.NotNil(t, result.Networking.VpcAccess)
	assert.Equal(t, "connector1", result.Networking.VpcAccess.Connector)
	assert.Equal(t, "ALL_TRAFFIC", result.Networking.VpcAccess.Egress)

	// Security
	assert.Equal(t, "sa@example.com", result.Security.ServiceAccount)
	assert.Equal(t, "key1", result.Security.EncryptionKey)
	assert.Equal(t, "default", result.Security.BinaryAuthorization)
	assert.Equal(t, "emergency", result.Security.BreakglassJustification)
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	ListServicesFunc  func(ctx context.Context, project, region string) ([]*runpb.Service, error)
	GetServiceFunc    func(ctx context.Context, name string) (*runpb.Service, error)
	UpdateServiceFunc func(ctx context.Context, service *runpb.Service) (*runpb.Service, error)
}

func (m *MockClient) ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error) {
	if m.ListServicesFunc != nil {
		return m.ListServicesFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) GetService(ctx context.Context, name string) (*runpb.Service, error) {
	if m.GetServiceFunc != nil {
		return m.GetServiceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
	if m.UpdateServiceFunc != nil {
		return m.UpdateServiceFunc(ctx, service)
	}
	return nil, nil
}

func TestList(t *testing.T) {
	// Save original client and restore after test
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		return []*runpb.Service{
			{Name: "projects/p/locations/r/services/s1", Uri: "uri1"},
			{Name: "projects/p/locations/r/services/s2", Uri: "uri2"},
		}, nil
	}

	services, err := List("p", "r")

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "s1", services[0].Name)
	assert.Equal(t, "s2", services[1].Name)
}

func TestUpdateScaling(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	// Setup GetService mock
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return &runpb.Service{
			Name: name,
			Scaling: &runpb.ServiceScaling{
				MinInstanceCount: 0,
				MaxInstanceCount: 1,
			},
		}, nil
	}

	// Setup UpdateService mock
	mock.UpdateServiceFunc = func(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
		// Assert that scaling was updated correctly
		assert.Equal(t, int32(2), service.Scaling.MinInstanceCount)
		assert.Equal(t, int32(5), service.Scaling.MaxInstanceCount)
		return service, nil
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "s1", 2, 5, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "s1", result.Name)
	assert.Equal(t, int32(2), result.Scaling.MinInstances)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		return nil, assert.AnError
	}

	services, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, services)
}

func TestUpdateScaling_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	// Test GetService Error
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return nil, assert.AnError
	}
	_, err := UpdateScaling(context.Background(), "p", "r", "s", 1, 2, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get service")

	// Test UpdateService Error
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return &runpb.Service{Name: name}, nil
	}
	mock.UpdateServiceFunc = func(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
		return nil, assert.AnError
	}
	_, err = UpdateScaling(context.Background(), "p", "r", "s", 1, 2, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update service")
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		if region == "us-central1" {
			return []*runpb.Service{
				{Name: "projects/p/locations/us-central1/services/s1"},
			}, nil
		}
		return []*runpb.Service{}, nil
	}

	services, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	
	// We expect at least 1 service
	found := false
	for _, s := range services {
		if s.Name == "s1" && s.Region == "us-central1" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find service s1 in us-central1")
}
