package create

import (
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
