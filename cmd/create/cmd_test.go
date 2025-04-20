package create

import (
	"bytes"
	"fmt"
	deps "h3s/cmd/dependencies"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/k3s"
	"h3s/internal/utils/cli"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommands(t *testing.T) {
	// Test that the main command is properly initialized
	assert.NotNil(t, Cmd)

	// Verify command structure - check subcommands
	var configCmd, credentialsCmd, clusterCmd *cobra.Command
	for _, cmd := range Cmd.Commands() {
		switch cmd.Use {
		case "config":
			configCmd = cmd
		case "credentials":
			credentialsCmd = cmd
		case "cluster":
			clusterCmd = cmd
		}
	}

	assert.NotNil(t, configCmd, "config subcommand should exist")
	assert.NotNil(t, credentialsCmd, "credentials subcommand should exist")
	assert.NotNil(t, clusterCmd, "cluster subcommand should exist")

	assert.Equal(t, "config", configCmd.Use)
	assert.Equal(t, "credentials", credentialsCmd.Use)
	assert.Equal(t, "cluster", clusterCmd.Use)
}

func TestRunCreateConfig(t *testing.T) {
	mockDeps := &deps.MockDependencies{
		MockGetK3sReleases: func(bool, bool, int) ([]k3s.Release, error) {
			return []k3s.Release{{Name: "v1.21.0"}}, nil
		},
		MockBuildClusterConfig: func(_ []k3s.Release) error {
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
			return &cluster.Cluster{
				Config: &config.Config{
					Domain: "test.domain",
					ControlPlane: config.ControlPlane{
						Pool: config.NodePool{
							Location: "test-location",
							Nodes:    1,
						},
					},
					WorkerPools: []config.NodePool{},
					K3sVersion:  "v1.21.0",
				},
			}, nil
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
		MockBuildClusterConfig: func(_ []k3s.Release) error {
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
			return &cluster.Cluster{
				Config: &config.Config{
					Domain: "test.domain",
					ControlPlane: config.ControlPlane{
						Pool: config.NodePool{
							Location: "test-location",
							Nodes:    1,
						},
					},
					WorkerPools: []config.NodePool{},
				},
			}, nil
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

// Helper to get a fresh create command for each test
func newCreateCmd() *cobra.Command {
	// Create a new instance of the command
	// We need to recreate the command structure as in cmd.go
	configCmd := cli.NewCommand(configConfig)
	credentialsCmd := cli.NewCommand(credentialsConfig)
	clusterCmd := cli.NewCommand(clusterConfig)

	// Create main command with subcommands
	createConfig.Subcommands = []*cobra.Command{configCmd, credentialsCmd, clusterCmd}
	return cli.NewCommand(createConfig)
}

func TestCreateCommand_FlagParsing(t *testing.T) {
	buf := new(bytes.Buffer)
	createCmd := newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)

	// Valid: create config --help
	createCmd.SetArgs([]string{"config", "--help"})
	err := createCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Usage:")
	buf.Reset()

	// Invalid: unknown flag
	createCmd = newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)
	createCmd.SetArgs([]string{"config", "--notaflag"})
	err = createCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "unknown flag: --notaflag")
	buf.Reset()

	// Invalid: extra argument
	createCmd = newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)
	createCmd.SetArgs([]string{"config", "extraarg"})
	err = createCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String()+err.Error(), "unknown command \"extraarg\"")
	buf.Reset()

	// Invalid: extra argument
	createCmd = newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)
	createCmd.SetArgs([]string{"config", "extraarg"})
	err = createCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String()+err.Error(), "unknown command \"extraarg\"")
	buf.Reset()

}

func TestCreateCommand_MissingRequiredArgs(t *testing.T) {
	buf := new(bytes.Buffer)
	createCmd := newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)

	// Simulate missing required args for cluster (if any required)
	createCmd.SetArgs([]string{"cluster"})
	err := createCmd.Execute()
	// The actual error message will depend on implementation, but should not be nil
	assert.Error(t, err)
}

func TestCreateCommand_OutputFormatting(t *testing.T) {
	buf := new(bytes.Buffer)
	createCmd := newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)

	// Valid output: config help
	createCmd.SetArgs([]string{"config", "--help"})
	err := createCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Usage:")
	buf.Reset()

	// Error output: unknown subcommand
	createCmd = newCreateCmd()
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)
	createCmd.SetArgs([]string{"notacommand"})
	err = createCmd.Execute()
	assert.Error(t, err)
	output = buf.String()
	assert.Contains(t, output, "unknown command \"notacommand\"")
}
