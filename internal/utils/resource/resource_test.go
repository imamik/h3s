package resource

import (
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

func TestResourceManager_Get(t *testing.T) {
	ctx := &cluster.Cluster{}
	manager := NewManager[*MockResource](ctx, logger.Server, "test-resource")

	// Test successful get
	resource, err := manager.Get(func() (*MockResource, error) {
		return &MockResource{ID: "123", Name: "test"}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID)
	assert.Equal(t, "test", resource.Name)

	// Test error case
	expectedErr := errors.New("get error")
	_, err = manager.Get(func() (*MockResource, error) {
		return nil, expectedErr
	})
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// We no longer check for nil resources in the generic context
	// This test is now obsolete
}

func TestResourceManager_Create(t *testing.T) {
	ctx := &cluster.Cluster{}
	manager := NewManager[*MockResource](ctx, logger.Server, "test-resource")

	// Test successful create
	resource, err := manager.Create(func() (*MockResource, error) {
		return &MockResource{ID: "123", Name: "test"}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID)
	assert.Equal(t, "test", resource.Name)

	// Test error case
	expectedErr := errors.New("create error")
	_, err = manager.Create(func() (*MockResource, error) {
		return nil, expectedErr
	})
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// We no longer check for nil resources in the generic context
	// This test is now obsolete
}

func TestResourceManager_Delete(t *testing.T) {
	ctx := &cluster.Cluster{}
	manager := NewManager[*MockResource](ctx, logger.Server, "test-resource")

	// Test successful delete
	err := manager.Delete(func() error {
		return nil
	})
	assert.NoError(t, err)

	// Test error case
	expectedErr := errors.New("delete error")
	err = manager.Delete(func() error {
		return expectedErr
	})
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestResourceManager_GetOrCreate(t *testing.T) {
	ctx := &cluster.Cluster{}
	manager := NewManager[*MockResource](ctx, logger.Server, "test-resource")

	// Test get succeeds
	resource, err := manager.GetOrCreate(
		func() (*MockResource, error) {
			return &MockResource{ID: "123", Name: "test"}, nil
		},
		func() (*MockResource, error) {
			return &MockResource{ID: "456", Name: "created"}, nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, "123", resource.ID) // Should return the "get" result

	// Test get fails but create succeeds
	resource, err = manager.GetOrCreate(
		func() (*MockResource, error) {
			return nil, errors.New("resource is nil")
		},
		func() (*MockResource, error) {
			return &MockResource{ID: "456", Name: "created"}, nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, "456", resource.ID) // Should return the "create" result

	// Test both get and create fail
	_, err = manager.GetOrCreate(
		func() (*MockResource, error) {
			return nil, errors.New("get error")
		},
		func() (*MockResource, error) {
			return nil, errors.New("create error")
		},
	)
	assert.Error(t, err)
	assert.Equal(t, "get error", err.Error()) // Should return the "get" error
}
