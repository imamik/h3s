package pool

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type MockCloudClient struct {
	// Add methods as needed for mocking
}

func (m *MockCloudClient) PlacementGroup() *hcloud.PlacementGroupClient {
	return &hcloud.PlacementGroupClient{}
}

func (m *MockCloudClient) Server() *hcloud.ServerClient {
	return &hcloud.ServerClient{}
}

// MockCluster is a minimal mock of cluster.Cluster for testing
// Only implements the methods and fields needed for CreatePools

type MockCluster struct {
	Config      *struct{ Name string }
	CloudClient *MockCloudClient
	Context     context.Context
}

func (c *MockCluster) GetName(names ...string) string {
	return "mock-cluster"
}

func (c *MockCluster) GetLabels(optionalLabels ...map[string]string) map[string]string {
	return map[string]string{"mock": "label"}
}
