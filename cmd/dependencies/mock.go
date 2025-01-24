package dependencies

import (
	"h3s/internal/cluster"
	"h3s/internal/k3s"
)

// MockDependencies provides a mock implementation for testing
type MockDependencies struct {
	MockGetClusterContext       func() (*cluster.Cluster, error)
	MockGetK3sReleases          func(bool, bool, int) ([]k3s.Release, error)
	MockInstallK3s              func(*cluster.Cluster) error
	MockCreateHetznerResources  func(*cluster.Cluster) error
	MockDestroyHetznerResources func(*cluster.Cluster) error
	MockInstallK8sComponents    func(*cluster.Cluster) error
	MockGenerateK8sToken        func(*cluster.Cluster, string, string, int) (string, error)
	MockDownloadKubeconfig      func(*cluster.Cluster) error
	MockExecuteSSHCommand       func(*cluster.Cluster, string) (string, error)
	MockExecuteLocalCommand     func(string) (string, error)
	MockBuildClusterConfig      func([]k3s.Release) error
	MockConfigureCredentials    func() error
	MockExecuteKubectlCommand   func(*cluster.Cluster, []string) (string, error)
	MockGetKubeconfigPath       func() (string, bool)
}

// GetClusterContext returns the cluster context using the mock implementation
func (m *MockDependencies) GetClusterContext() (*cluster.Cluster, error) {
	return m.MockGetClusterContext()
}

// GetK3sReleases implements the interface method
func (m *MockDependencies) GetK3sReleases(stable, latest bool, limit int) ([]k3s.Release, error) {
	return m.MockGetK3sReleases(stable, latest, limit)
}

// BuildClusterConfig implements the interface method
func (m *MockDependencies) BuildClusterConfig(releases []k3s.Release) error {
	return m.MockBuildClusterConfig(releases)
}

// InstallK3s implements the interface method
func (m *MockDependencies) InstallK3s(ctx *cluster.Cluster) error {
	return m.MockInstallK3s(ctx)
}

// CreateHetznerResources implements the interface method
func (m *MockDependencies) CreateHetznerResources(ctx *cluster.Cluster) error {
	return m.MockCreateHetznerResources(ctx)
}

// DestroyHetznerResources implements the interface method
func (m *MockDependencies) DestroyHetznerResources(ctx *cluster.Cluster) error {
	return m.MockDestroyHetznerResources(ctx)
}

// InstallK8sComponents implements the interface method
func (m *MockDependencies) InstallK8sComponents(ctx *cluster.Cluster) error {
	return m.MockInstallK8sComponents(ctx)
}

// GenerateK8sToken implements the interface method
func (m *MockDependencies) GenerateK8sToken(ctx *cluster.Cluster, namespace, serviceAccount string, hours int) (string, error) {
	return m.MockGenerateK8sToken(ctx, namespace, serviceAccount, hours)
}

// DownloadKubeconfig implements the interface method
func (m *MockDependencies) DownloadKubeconfig(ctx *cluster.Cluster) error {
	return m.MockDownloadKubeconfig(ctx)
}

// ExecuteSSHCommand implements the interface method
func (m *MockDependencies) ExecuteSSHCommand(ctx *cluster.Cluster, command string) (string, error) {
	return m.MockExecuteSSHCommand(ctx, command)
}

// ExecuteLocalCommand implements the interface method
func (m *MockDependencies) ExecuteLocalCommand(command string) (string, error) {
	return m.MockExecuteLocalCommand(command)
}

// ConfigureCredentials implements the interface method
func (m *MockDependencies) ConfigureCredentials() error {
	return m.MockConfigureCredentials()
}

// ExecuteKubectlCommand implements the interface method
func (m *MockDependencies) ExecuteKubectlCommand(ctx *cluster.Cluster, args []string) (string, error) {
	return m.MockExecuteKubectlCommand(ctx, args)
}

// GetKubeconfigPath implements the interface method
func (m *MockDependencies) GetKubeconfigPath() (string, bool) {
	return m.MockGetKubeconfigPath()
}
