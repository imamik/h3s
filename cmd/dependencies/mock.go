package dependencies

import (
	"h3s/internal/cluster"
)

// MockDependencies provides a mock implementation for testing
type MockDependencies struct {
	MockGetClusterContext      func() (*cluster.Cluster, error)
	MockGetK3sReleases         func(bool, bool, int) ([]string, error)
	MockInstallK3s             func(*cluster.Cluster) error
	MockCreateHetznerResources func(*cluster.Cluster) error
	// Add other mock method signatures as needed
}

// Implement each method with the corresponding mock function
func (m *MockDependencies) GetClusterContext() (*cluster.Cluster, error) {
	return m.MockGetClusterContext()
}

// Implement other interface methods similarly...
