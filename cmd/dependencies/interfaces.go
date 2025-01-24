package dependencies

import (
	"h3s/internal/cluster"
	"h3s/internal/k3s"
)

// CommandDependencies defines the contract for CLI command dependencies
type CommandDependencies interface {
	// Cluster Context and Management
	GetClusterContext() (*cluster.Cluster, error)

	// K3s Related
	GetK3sReleases(filterStable, filterLatest bool, limit int) ([]k3s.Release, error)
	InstallK3s(ctx *cluster.Cluster) error

	// Hetzner Cloud Resources
	CreateHetznerResources(ctx *cluster.Cluster) error
	DestroyHetznerResources(ctx *cluster.Cluster) error

	// Kubernetes Components
	InstallK8sComponents(ctx *cluster.Cluster) error
	GenerateK8sToken(ctx *cluster.Cluster, namespace, serviceAccount string, hours int) (string, error)

	// Kubeconfig Management
	DownloadKubeconfig(ctx *cluster.Cluster) error

	// SSH and Command Execution
	ExecuteSSHCommand(ctx *cluster.Cluster, command string) (string, error)
	ExecuteLocalCommand(command string) (string, error)

	// Configuration Management
	BuildClusterConfig(k3sReleases []k3s.Release) error
	ConfigureCredentials() error

	// Kubectl Related
	ExecuteKubectlCommand(ctx *cluster.Cluster, args []string) (string, error)
	GetKubeconfigPath() (string, bool)
}
