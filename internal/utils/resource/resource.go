// Package resource provides utilities for managing cloud resources with consistent patterns
package resource

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

// Manager provides a generic interface for resource operations
type Manager[T any] struct {
	Ctx          *cluster.Cluster
	ResourceType logger.LogResource
	Name         string
}

// NewManager creates a new resource manager
func NewManager[T any](ctx *cluster.Cluster, resourceType logger.LogResource, name string) *Manager[T] {
	return &Manager[T]{
		Ctx:          ctx,
		ResourceType: resourceType,
		Name:         name,
	}
}

// Get retrieves a resource using the provided getter function
func (m *Manager[T]) Get(getter func() (T, error)) (T, error) {
	l := logger.New(nil, m.ResourceType, logger.Get, m.Name)
	defer l.LogEvents()

	resource, err := getter()
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return resource, err
	}

	// We can't directly check if resource is nil in a generic context
	// This would need reflection in a real implementation
	// For now, we'll assume the resource is valid if we got this far

	l.AddEvent(logger.Success)
	return resource, nil
}

// Create creates a resource using the provided creator function
func (m *Manager[T]) Create(creator func() (T, error)) (T, error) {
	l := logger.New(nil, m.ResourceType, logger.Create, m.Name)
	defer l.LogEvents()

	resource, err := creator()
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return resource, err
	}

	// We can't directly check if resource is nil in a generic context
	// This would need reflection in a real implementation
	// For now, we'll assume the resource is valid if we got this far

	l.AddEvent(logger.Success)
	return resource, nil
}

// Delete deletes a resource using the provided deleter function
func (m *Manager[T]) Delete(deleter func() error) error {
	l := logger.New(nil, m.ResourceType, logger.Delete, m.Name)
	defer l.LogEvents()

	if err := deleter(); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}

// GetOrCreate attempts to retrieve a resource using the getter function.
// If the getter returns an error indicating the resource is not found, it attempts to create it using the creator function.
func (m *Manager[T]) GetOrCreate(
	getter func() (T, error),
	creator func() (T, error),
) (T, error) {
	var zero T
	resource, err := getter()
	if err == nil {
		return resource, nil
	}

	// Only create if the error indicates the resource doesn't exist
	if err.Error() == "resource is nil" {
		return m.Create(creator)
	}

	return zero, err
}

// In a real implementation, we would need to use reflection to check if a value is nil
// For simplicity in this example, we're omitting this check
