// Package naming provides utilities for generating consistent resource names
package naming

import (
	"h3s/internal/cluster"
	"strings"
)

// ResourceName generates a consistent name for a resource based on the cluster context and components
func ResourceName(ctx *cluster.Cluster, components ...string) string {
	if ctx == nil || ctx.Config == nil {
		return "invalid-cluster"
	}
	return ctx.GetName(components...)
}

// FormatName joins name components with a consistent separator
func FormatName(prefix string, components ...string) string {
	allComponents := append([]string{prefix}, components...)
	return strings.Join(allComponents, "-")
}
