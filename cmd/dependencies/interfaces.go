package dependencies

import (
	"h3s/internal/cluster"
	"h3s/internal/k3s"
)

// CommandDependencies defines the contract for all dependencies required by CLI commands.
//
// CONTRACT REQUIREMENTS:
//   - All methods must be implemented by any struct claiming to satisfy this interface.
//   - Implementations must not panic; errors should be returned using the error return value.
//   - Implementations should be concurrency-safe if used in concurrent CLI commands.
//   - For resource creation/deletion (e.g., CreateHetznerResources), implementations must ensure idempotency.
//   - For context-related methods (e.g., GetClusterContext), implementations must return fully initialized, valid objects or a descriptive error.
//   - Implementations are responsible for handling provider-specific errors and surfacing actionable messages.
//   - All returned objects must meet the expected types and invariants described in method comments.
//
// IMPLEMENTATION RESPONSIBILITIES:
//   - Document any side effects (e.g., network calls, file system writes).
//   - Ensure all configuration and credential requirements are met before performing any action.
//   - Log relevant events for debugging and traceability.
//   - Adhere to the security guidelines of the project (do not log secrets, etc.).
//
// CONTRACT ENFORCEMENT:
//   - All implementations must be registered in the contract test (see contract_test.go).
//   - The contract test will fail if a method is missing or has the wrong signature.
//   - For runtime contract checks, see cmd/dependencies/contract_test.go.
//
// See also: https://golang.org/doc/effective_go#interfaces

type CommandDependencies interface {
	// Cluster Context and Management
	// GetClusterContext returns the current cluster context.
	// The returned cluster object must be fully initialized and valid.
	// If an error occurs, a descriptive error message should be returned.
	GetClusterContext() (*cluster.Cluster, error)

	// K3s Related
	// GetK3sReleases returns a list of K3s releases filtered by stability and latest version.
	// The filterStable parameter determines whether to include only stable releases.
	// The filterLatest parameter determines whether to include only the latest release.
	// The limit parameter determines the maximum number of releases to return.
	// The returned releases must meet the expected types and invariants described in the k3s.Release type.
	GetK3sReleases(filterStable, filterLatest bool, limit int) ([]k3s.Release, error)
	// InstallK3s installs K3s on the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// If an error occurs, a descriptive error message should be returned.
	InstallK3s(ctx *cluster.Cluster) error

	// Hetzner Cloud Resources
	// CreateHetznerResources creates Hetzner cloud resources for the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// Implementations must ensure idempotency.
	// If an error occurs, a descriptive error message should be returned.
	CreateHetznerResources(ctx *cluster.Cluster) error
	// DestroyHetznerResources destroys Hetzner cloud resources for the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// Implementations must ensure idempotency.
	// If an error occurs, a descriptive error message should be returned.
	DestroyHetznerResources(ctx *cluster.Cluster) error

	// Kubernetes Components
	// InstallK8sComponents installs Kubernetes components on the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// If an error occurs, a descriptive error message should be returned.
	InstallK8sComponents(ctx *cluster.Cluster) error
	// GenerateK8sToken generates a Kubernetes token for the specified namespace and service account.
	// The ctx parameter must be a valid cluster context.
	// The namespace and serviceAccount parameters must be valid.
	// The hours parameter determines the token's validity period.
	// The returned token must meet the expected types and invariants described in the string type.
	GenerateK8sToken(ctx *cluster.Cluster, namespace, serviceAccount string, hours int) (string, error)

	// Kubeconfig Management
	// DownloadKubeconfig downloads the kubeconfig file for the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// If an error occurs, a descriptive error message should be returned.
	DownloadKubeconfig(ctx *cluster.Cluster) error

	// SSH and Command Execution
	// ExecuteSSHCommand executes an SSH command on the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// The command parameter must be a valid SSH command.
	// The returned output must meet the expected types and invariants described in the string type.
	ExecuteSSHCommand(ctx *cluster.Cluster, command string) (string, error)
	// ExecuteLocalCommand executes a local command.
	// The command parameter must be a valid local command.
	// The returned output must meet the expected types and invariants described in the string type.
	ExecuteLocalCommand(command string) (string, error)

	// Configuration Management
	// BuildClusterConfig builds the cluster configuration from the specified K3s releases.
	// The k3sReleases parameter must be a valid list of K3s releases.
	// If an error occurs, a descriptive error message should be returned.
	BuildClusterConfig(k3sReleases []k3s.Release) error
	// ConfigureCredentials configures the credentials for the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// If an error occurs, a descriptive error message should be returned.
	ConfigureCredentials() error

	// Kubectl Related
	// ExecuteKubectlCommand executes a kubectl command on the specified cluster context.
	// The ctx parameter must be a valid cluster context.
	// The args parameter must be a valid list of kubectl command arguments.
	// The returned output must meet the expected types and invariants described in the string type.
	ExecuteKubectlCommand(ctx *cluster.Cluster, args []string) (string, error)
	// GetKubeconfigPath returns the path to the kubeconfig file.
	// The returned path must meet the expected types and invariants described in the string type.
	GetKubeconfigPath() (string, bool)
}
