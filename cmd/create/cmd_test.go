package create

import (
	"fmt"
	deps "h3s/cmd/dependencies"
	"h3s/internal/cluster"
	"h3s/internal/k3s"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommands(t *testing.T) {
	// Test that all subcommands are properly initialized
	assert.NotNil(t, createConfigCmd)
	assert.NotNil(t, createCredentialsCmd)
	assert.NotNil(t, createClusterCmd)

	// Verify command structure
	assert.Equal(t, "config", createConfigCmd.Use)
	assert.Equal(t, "credentials", createCredentialsCmd.Use)
	assert.Equal(t, "cluster", createClusterCmd.Use)
}

func TestRunCreateConfig(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetK3sReleases: func(bool, bool, int) ([]k3s.Release, error) {
			return []k3s.Release{{Name: "v1.21.0"}}, nil
		},
		MockBuildClusterConfig: func(releases []k3s.Release) error {
			return nil
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateConfig(nil, nil)
	assert.NoError(t, err)
}

func TestRunCreateCluster(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetClusterContext: func() (*cluster.Cluster, error) {
			return &cluster.Cluster{}, nil
		},
		MockCreateHetznerResources: func(*cluster.Cluster) error {
			return nil
		},
		MockInstallK3s: func(*cluster.Cluster) error {
			return nil
		},
		MockInstallK8sComponents: func(*cluster.Cluster) error {
			return nil
		},
		MockDownloadKubeconfig: func(*cluster.Cluster) error {
			return nil
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateCluster(nil, nil)
	assert.NoError(t, err)
}

func TestRunCreateConfig_K3sReleasesError(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetK3sReleases: func(bool, bool, int) ([]k3s.Release, error) {
			return nil, fmt.Errorf("failed to fetch releases")
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateConfig(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get k3s releases")
}

func TestRunCreateConfig_BuildConfigError(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetK3sReleases: func(bool, bool, int) ([]k3s.Release, error) {
			return []k3s.Release{{Name: "v1.21.0"}}, nil
		},
		MockBuildClusterConfig: func(releases []k3s.Release) error {
			return fmt.Errorf("failed to build config")
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateConfig(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build configuration")
}

func TestRunCreateCredentials_Error(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockConfigureCredentials: func() error {
			return fmt.Errorf("failed to configure credentials")
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateCredentials(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to configure credentials")
}

func TestRunCreateCluster_GetContextError(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetClusterContext: func() (*cluster.Cluster, error) {
			return nil, fmt.Errorf("failed to get context")
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateCluster(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load cluster context")
}

func TestRunCreateCluster_CreateHetznerError(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetClusterContext: func() (*cluster.Cluster, error) {
			return &cluster.Cluster{}, nil
		},
		MockCreateHetznerResources: func(*cluster.Cluster) error {
			return fmt.Errorf("failed to create resources")
		},
	}
	deps.Get = func() deps.CommandDependencies {
		return mockDeps
	}

	err := runCreateCluster(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create hetzner resources")
}
