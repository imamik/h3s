// Package testutils provides utilities for testing CLI commands
package testutils

import (
	"bytes"
	"h3s/cmd/dependencies"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/k3s"
	"time"

	"github.com/spf13/cobra"
)

// SetupTestCmd sets up a command for testing, swapping dependencies for mocks and capturing output.
// Returns a buffer for capturing output and a cleanup function.
func SetupTestCmd(cmd *cobra.Command, mockDeps *dependencies.MockDependencies) (*bytes.Buffer, func()) {
	buf := new(bytes.Buffer)
	origOut := cmd.OutOrStdout()
	cmd.SetOut(buf)

	// Swap dependencies to mock
	origDeps := dependencies.Get
	dependencies.Get = func() dependencies.CommandDependencies {
		return mockDeps
	}

	cleanup := func() {
		cmd.SetOut(origOut)
		dependencies.Get = origDeps
	}

	return buf, cleanup
}

// CreateMockCluster creates a mock cluster for testing
func CreateMockCluster() *cluster.Cluster {
	return &cluster.Cluster{
		Config: &config.Config{
			Domain: "test.domain",
			ControlPlane: config.ControlPlane{
				Pool: config.NodePool{
					Nodes:    1,
					Name:     "control-plane",
					Instance: "cx21",
					Location: "nbg1",
				},
				AsWorkerPool: false,
			},
			K3sVersion:  "v1.28.2+k3s1",
			Name:        "test-cluster",
			NetworkZone: "eu-central",
			SSHKeyPaths: config.SSHKeyPaths{
				PrivateKeyPath: "/tmp/id_rsa",
				PublicKeyPath:  "/tmp/id_rsa.pub",
			},
			CertManager: config.CertManager{
				Email:      "test@example.com",
				Production: false,
			},
			WorkerPools: []config.NodePool{
				{
					Nodes:    1,
					Name:     "worker",
					Instance: "cx21",
					Location: "nbg1",
				},
			},
		},
		Credentials: &credentials.ProjectCredentials{
			K3sToken:        "test-token",
			HCloudToken:     "test-hcloud-token",
			HetznerDNSToken: "test-dns-token",
		},
	}
}

// DefaultMockDependencies returns a MockDependencies with default implementations
func DefaultMockDependencies() *dependencies.MockDependencies {
	mockCluster := CreateMockCluster()

	return &dependencies.MockDependencies{
		MockGetClusterContext: func() (*cluster.Cluster, error) {
			return mockCluster, nil
		},
		MockGetK3sReleases: func(_, _ bool, _ int) ([]k3s.Release, error) {
			return []k3s.Release{
				{
					Name:        "v1.28.2+k3s1",
					PublishedAt: time.Now().Add(-24 * time.Hour),
					Prerelease:  false,
					Draft:       false,
				},
			}, nil
		},
		MockInstallK3s: func(_ *cluster.Cluster) error {
			return nil
		},
		MockInstallK8sComponents: func(_ *cluster.Cluster) error {
			return nil
		},
		MockExecuteSSHCommand: func(_ *cluster.Cluster, command string) (string, error) {
			return "mock ssh output: " + command, nil
		},
		MockExecuteLocalCommand: func(command string) (string, error) {
			return "mock local output: " + command, nil
		},
		MockGetKubeconfigPath: func() (string, bool) {
			return "mock-kubeconfig.yaml", true
		},
	}
}
