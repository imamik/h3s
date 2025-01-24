// Package dependencies provides the dependencies for the commands
package dependencies

import (
	"h3s/internal/cluster"
	"h3s/internal/config/create"
	"h3s/internal/config/credentials"
	"h3s/internal/config/path"
	"h3s/internal/hetzner"
	"h3s/internal/k3s"
	"h3s/internal/k8s"
	"h3s/internal/k8s/kubeconfig"
	"h3s/internal/utils/common"
	"h3s/internal/utils/execute"
	"h3s/internal/utils/file"
	"h3s/internal/utils/kubectl"
)

// DefaultDependencies provides the standard implementation of CommandDependencies
type DefaultDependencies struct{}

// GetClusterContext returns the cluster context
func (d *DefaultDependencies) GetClusterContext() (*cluster.Cluster, error) {
	return cluster.Context()
}

// GetK3sReleases returns the k3s releases
func (d *DefaultDependencies) GetK3sReleases(filterStable, filterLatest bool, limit int) ([]k3s.Release, error) {
	return k3s.GetFilteredReleases(filterStable, filterLatest, limit)
}

// InstallK3s installs k3s to the nodes in the cluster
func (d *DefaultDependencies) InstallK3s(ctx *cluster.Cluster) error {
	return k3s.Install(ctx)
}

// CreateHetznerResources creates the hetzner resources for the cluster
func (d *DefaultDependencies) CreateHetznerResources(ctx *cluster.Cluster) error {
	return hetzner.Create(ctx)
}

// DestroyHetznerResources destroys the hetzner resources for the cluster
func (d *DefaultDependencies) DestroyHetznerResources(ctx *cluster.Cluster) error {
	return hetzner.Destroy(ctx)
}

// InstallK8sComponents installs the k8s components to the main node in the cluster
func (d *DefaultDependencies) InstallK8sComponents(ctx *cluster.Cluster) error {
	return k8s.Install(ctx)
}

// GenerateK8sToken generates a k8s token for the given namespace, service account, and hours
func (d *DefaultDependencies) GenerateK8sToken(ctx *cluster.Cluster, namespace, serviceAccount string, hours int) (string, error) {
	return k8s.Token(ctx, namespace, serviceAccount, hours)
}

// DownloadKubeconfig downloads the kubeconfig from the cluster
func (d *DefaultDependencies) DownloadKubeconfig(ctx *cluster.Cluster) error {
	return kubeconfig.Download(ctx)
}

// ExecuteSSHCommand executes an ssh command on the cluster
func (d *DefaultDependencies) ExecuteSSHCommand(ctx *cluster.Cluster, command string) (string, error) {
	return common.SSH(ctx, command)
}

// ExecuteLocalCommand executes a local command
func (d *DefaultDependencies) ExecuteLocalCommand(command string) (string, error) {
	return execute.Local(command)
}

// BuildClusterConfig builds the cluster config, e.g. the amount of nodes, the dns name, etc.
func (d *DefaultDependencies) BuildClusterConfig(k3sReleases []k3s.Release) error {
	return create.Build(k3sReleases)
}

// ConfigureCredentials configures the credentials, e.g. hetzner api key & dns token, etc.
func (d *DefaultDependencies) ConfigureCredentials() error {
	return credentials.Configure()
}

// ExecuteKubectlCommand executes a kubectl command either via kubeconfig or SSH
func (d *DefaultDependencies) ExecuteKubectlCommand(ctx *cluster.Cluster, args []string) (string, error) {
	kubeconfigPath, exists := d.GetKubeconfigPath()
	if exists {
		cmd, err := kubectl.New(args...).AddKubeConfigPath(kubeconfigPath).String()
		if err != nil {
			return "", err
		}
		return execute.Local(cmd)
	}

	cmd, err := kubectl.New(args...).EmbedFileContent().String()
	if err != nil {
		return "", err
	}
	return common.SSH(ctx, cmd)
}

// GetKubeconfigPath returns the path to kubeconfig if it exists
func (d *DefaultDependencies) GetKubeconfigPath() (string, bool) {
	p := string(path.KubeConfigFileName)
	f := file.New(p)
	absPath, err := f.Path()
	if err != nil {
		return "", false
	}
	return absPath, f.Exists()
}

// NewDefaultDependencies creates a new instance of DefaultDependencies
func NewDefaultDependencies() CommandDependencies {
	return &DefaultDependencies{}
}
