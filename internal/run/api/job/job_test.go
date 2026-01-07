package job

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	ListJobsFunc func(ctx context.Context, project, region string) ([]*runpb.Job, error)
	RunJobFunc   func(ctx context.Context, name string) (*runpb.Execution, error)
}

func (m *MockClient) ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error) {
	if m.ListJobsFunc != nil {
		return m.ListJobsFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) RunJob(ctx context.Context, name string) (*runpb.Execution, error) {
	if m.RunJobFunc != nil {
		return m.RunJobFunc(ctx, name)
	}
	return nil, nil
}

func TestMapJob(t *testing.T) {
	// ... (existing test)
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

func TestMapJob_NilFields(t *testing.T) {
	resp := &runpb.Job{
		Name:    "projects/my-project/locations/us-central1/jobs/my-job",
		Creator: "user@example.com",
	}

	result := mapJob(resp, "us-central1")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Nil(t, result.LatestCreatedExecution)
	assert.Nil(t, result.TerminalCondition)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		return []*runpb.Job{
			{Name: "job1"},
			{Name: "job2"},
		}, nil
	}

	jobs, err := List("p", "r")
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, "job1", jobs[0].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		return nil, assert.AnError
	}

	jobs, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, jobs)
}

func TestExecute(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.RunJobFunc = func(ctx context.Context, name string) (*runpb.Execution, error) {
		assert.Equal(t, "projects/p/locations/r/jobs/myjob", name)
		return &runpb.Execution{Name: "exec1"}, nil
	}

	exec, err := Execute("p", "r", "myjob")
	assert.NoError(t, err)
	assert.NotNil(t, exec)
	assert.Equal(t, "exec1", exec.Name)
}

func TestExecute_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.RunJobFunc = func(ctx context.Context, name string) (*runpb.Execution, error) {
		return nil, assert.AnError
	}

	exec, err := Execute("p", "r", "myjob")
	assert.Error(t, err)
	assert.Nil(t, exec)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		if region == "us-central1" {
			return []*runpb.Job{{Name: "job-us"}}, nil
		}
		return []*runpb.Job{}, nil
	}

	jobs, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	
	found := false
	for _, j := range jobs {
		if j.Name == "job-us" && j.Region == "us-central1" {
			found = true
			break
		}
	}
	assert.True(t, found)
}
