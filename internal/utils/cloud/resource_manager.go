// Package cloud provides utilities for interacting with cloud providers
package cloud

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/resource"
)

// createResourceManager is a helper function to create a resource manager
func createResourceManager[T any](ctx *cluster.Cluster, resourceType logger.LogResource, name string) *resource.Manager[T] {
	return resource.NewManager[T](ctx, resourceType, name)
}
