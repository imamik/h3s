// Package path provides types and constants for file names used in the application.
package path

// FileName is a type alias for a string representing a file name.
type FileName string

const ( // #nosec G101 -- These are configuration file names, not credentials
	SecretsFileName    FileName = "h3s-secrets.yaml"    // SecretsFileName is the file name for the secrets configuration file containing e.g. Hetzner cloud tokens
	KubeConfigFileName FileName = "h3s-kubeconfig.yaml" // KubeConfigFileName is the file name for the h3s cluster configuration file
)
