package workerpool

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient
type MockClient struct {
	ListWorkerPoolsFunc  func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error)
	GetWorkerPoolFunc    func(ctx context.Context, name string) (*runpb.WorkerPool, error)
	UpdateWorkerPoolFunc func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error)
}

func (m *MockClient) ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
	if m.ListWorkerPoolsFunc != nil {
		return m.ListWorkerPoolsFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error) {
	if m.GetWorkerPoolFunc != nil {
		return m.GetWorkerPoolFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
	if m.UpdateWorkerPoolFunc != nil {
		return m.UpdateWorkerPoolFunc(ctx, workerPool)
	}
	return nil, nil
}

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

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return []*runpb.WorkerPool{{Name: "pool1"}}, nil
	}

	pools, err := List("p", "r")
	assert.NoError(t, err)
	assert.Len(t, pools, 1)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	pools, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, pools)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return []*runpb.WorkerPool{{Name: "projects/" + project + "/locations/" + region + "/workerPools/pool-" + region}}, nil
	}

	pools, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	// We expect at least one pool per region. api_region.List() has > 0 regions.
	assert.NotEmpty(t, pools)
}

func TestUpdateScaling(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return &runpb.WorkerPool{Name: name}, nil
	}
	mock.UpdateWorkerPoolFunc = func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
		assert.Equal(t, int32(5), *workerPool.Scaling.ManualInstanceCount)
		return workerPool, nil
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int32(5), result.Scaling.ManualInstanceCount)
}

func TestUpdateScaling_GetError(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get worker pool")
}

func TestUpdateScaling_UpdateError(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return &runpb.WorkerPool{Name: name}, nil
	}
	mock.UpdateWorkerPoolFunc = func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update worker pool")
}
