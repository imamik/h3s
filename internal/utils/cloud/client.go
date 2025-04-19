// Package cloud provides utilities for interacting with cloud providers
package cloud

import (
	"context"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

// ResourceClient provides a generic interface for cloud resource operations
type ResourceClient[T any] struct {
	Ctx          *cluster.Cluster
	ResourceType logger.LogResource
	ResourceName string
}

// NewClient creates a new cloud resource client
func NewClient[T any](ctx *cluster.Cluster, resourceType logger.LogResource, resourceName string) *ResourceClient[T] {
	return &ResourceClient[T]{
		Ctx:          ctx,
		ResourceType: resourceType,
		ResourceName: resourceName,
	}
}

// Get retrieves a resource using the provided getter function
func (c *ResourceClient[T]) Get(getter func(context.Context, string) (T, error)) (T, error) {
	resourceManager := createResourceManager[T](c.Ctx, c.ResourceType, c.ResourceName)
	return resourceManager.Get(func() (T, error) {
		return getter(c.Ctx.Context, c.ResourceName)
	})
}

// Create creates a resource using the provided creator function
func (c *ResourceClient[T]) Create(creator func(context.Context, interface{}) (T, error), opts interface{}) (T, error) {
	resourceManager := createResourceManager[T](c.Ctx, c.ResourceType, c.ResourceName)
	return resourceManager.Create(func() (T, error) {
		return creator(c.Ctx.Context, opts)
	})
}

// Delete deletes a resource using the provided deleter function
func (c *ResourceClient[T]) Delete(deleter func(context.Context, T) error, resource T) error {
	resourceManager := createResourceManager[T](c.Ctx, c.ResourceType, c.ResourceName)
	return resourceManager.Delete(func() error {
		return deleter(c.Ctx.Context, resource)
	})
}

// GetOrCreate gets a resource if it exists, or creates it if it doesn't
func (c *ResourceClient[T]) GetOrCreate(
	getter func(context.Context, string) (T, error),
	creator func(context.Context, interface{}) (T, error),
	opts interface{},
) (T, error) {
	resourceManager := createResourceManager[T](c.Ctx, c.ResourceType, c.ResourceName)
	return resourceManager.GetOrCreate(
		func() (T, error) {
			return getter(c.Ctx.Context, c.ResourceName)
		},
		func() (T, error) {
			return creator(c.Ctx.Context, opts)
		},
	)
}
