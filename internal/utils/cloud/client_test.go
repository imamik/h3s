package cloud

import (
	"context"
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock resource for testing
type MockResource struct {
	ID   string
	Name string
}

// Mock options for testing
type MockOptions struct {
	Name string
}

// MockCluster is a minimal mock of cluster.Cluster for testing
type MockCluster struct {
	Context context.Context
}

func TestResourceClient_Get(t *testing.T) {
	// Create a mock cluster
	ctx := &cluster.Cluster{
		Context: context.Background(),
	}

	// Create a client
	client := NewClient[*MockResource](ctx, logger.Server, "test-resource")

	// Test successful get
	resource, err := client.Get(func(_ context.Context, name string) (*MockResource, error) {
		return &MockResource{ID: "123", Name: name}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID)
	assert.Equal(t, "test-resource", resource.Name)

	// Test error case
	expectedErr := errors.New("get error")
	_, err = client.Get(func(_ context.Context, _ string) (*MockResource, error) {
		return nil, expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestResourceClient_Create(t *testing.T) {
	// Create a mock cluster
	ctx := &cluster.Cluster{
		Context: context.Background(),
	}

	// Create a client
	client := NewClient[*MockResource](ctx, logger.Server, "test-resource")

	// Test successful create
	opts := &MockOptions{Name: "test"}
	resource, err := client.Create(func(_ context.Context, opts interface{}) (*MockResource, error) {
		options := opts.(*MockOptions)
		return &MockResource{ID: "123", Name: options.Name}, nil
	}, opts)

	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID)
	assert.Equal(t, "test", resource.Name)

	// Test error case
	expectedErr := errors.New("create error")
	_, err = client.Create(func(_ context.Context, _ interface{}) (*MockResource, error) {
		return nil, expectedErr
	}, opts)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestResourceClient_Delete(t *testing.T) {
	// Create a mock cluster
	ctx := &cluster.Cluster{
		Context: context.Background(),
	}

	// Create a client
	client := NewClient[*MockResource](ctx, logger.Server, "test-resource")

	// Create a resource to delete
	resource := &MockResource{ID: "123", Name: "test"}

	// Test successful delete
	err := client.Delete(func(_ context.Context, _ *MockResource) error {
		return nil
	}, resource)

	assert.NoError(t, err)

	// Test error case
	expectedErr := errors.New("delete error")
	err = client.Delete(func(_ context.Context, _ *MockResource) error {
		return expectedErr
	}, resource)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestResourceClient_GetOrCreate(t *testing.T) {
	// Create a mock cluster
	ctx := &cluster.Cluster{
		Context: context.Background(),
	}

	// Create a client
	client := NewClient[*MockResource](ctx, logger.Server, "test-resource")

	// Test get succeeds
	opts := &MockOptions{Name: "test"}
	resource, err := client.GetOrCreate(
		func(_ context.Context, name string) (*MockResource, error) {
			return &MockResource{ID: "123", Name: name}, nil
		},
		func(_ context.Context, opts interface{}) (*MockResource, error) {
			options := opts.(*MockOptions)
			return &MockResource{ID: "456", Name: options.Name}, nil
		},
		opts,
	)

	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID) // Should return the "get" result

	// Test get fails but create succeeds
	resource, err = client.GetOrCreate(
		func(_ context.Context, _ string) (*MockResource, error) {
			return nil, errors.New("resource is nil")
		},
		func(_ context.Context, opts interface{}) (*MockResource, error) {
			options := opts.(*MockOptions)
			return &MockResource{ID: "456", Name: options.Name}, nil
		},
		opts,
	)

	assert.NoError(t, err)
	assert.Equal(t, "456", resource.ID) // Should return the "create" result

	// Test both get and create fail
	_, err = client.GetOrCreate(
		func(_ context.Context, _ string) (*MockResource, error) {
			return nil, errors.New("get error")
		},
		func(_ context.Context, _ interface{}) (*MockResource, error) {
			return nil, errors.New("create error")
		},
		opts,
	)

	assert.Error(t, err)
	assert.Equal(t, "get error", err.Error()) // Should return the "get" error
}
