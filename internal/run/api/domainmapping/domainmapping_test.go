package domainmapping

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/run/v1"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	ListDomainMappingsFunc func(ctx context.Context, project, region string) ([]*run.DomainMapping, error)
}

func (m *MockClient) ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
	if m.ListDomainMappingsFunc != nil {
		return m.ListDomainMappingsFunc(ctx, project, region)
	}
	return nil, nil
}

func TestMapDomainMapping(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	resp := &run.DomainMapping{
		Metadata: &run.ObjectMeta{
			Name:              "example.com",
			CreationTimestamp: now,
		},
		Spec: &run.DomainMappingSpec{
			RouteName: "my-service",
		},
		Status: &run.DomainMappingStatus{
			ResourceRecords: []*run.ResourceRecord{
				{
					Type:   "CNAME",
					Name:   "www",
					Rrdata: "ghs.googlehosted.com.",
				},
			},
			Conditions: []*run.GoogleCloudRunV1Condition{
				{
					Type:    "Ready",
					Status:  "True",
					Message: "Ready",
				},
			},
		},
	}

	result := mapDomainMapping(resp, "my-project", "us-central1")

	assert.Equal(t, "example.com", result.Name)
	assert.Equal(t, "my-service", result.RouteName)
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "us-central1", result.Region)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, "CNAME", result.Records[0].Type)
	assert.Equal(t, "www", result.Records[0].Name)
	assert.Equal(t, "ghs.googlehosted.com.", result.Records[0].RRData)
	assert.Len(t, result.Conditions, 1)
	assert.Equal(t, "Ready", result.Conditions[0].Type)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListDomainMappingsFunc = func(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
		return []*run.DomainMapping{
			{
				Metadata: &run.ObjectMeta{Name: "d1"},
				Spec:     &run.DomainMappingSpec{RouteName: "s1"},
			},
			{
				Metadata: &run.ObjectMeta{Name: "d2"},
				Spec:     &run.DomainMappingSpec{RouteName: "s2"},
			},
		}, nil
	}

	dms, err := List("p", "r")

	assert.NoError(t, err)
	assert.Len(t, dms, 2)
	assert.Equal(t, "d1", dms[0].Name)
	assert.Equal(t, "s1", dms[0].RouteName)
	assert.Equal(t, "d2", dms[1].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListDomainMappingsFunc = func(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
		return nil, assert.AnError
	}

	dms, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, dms)
}