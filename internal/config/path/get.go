// Package path provides types and constants for file names used in the application.
package path

import (
	"os"
)

// FileName is a type alias for a string representing a file name.
type FileName string

// SecretsFileName returns the filename for the secrets configuration file, honoring H3S_CREDENTIALS env var if set
func SecretsFileName() FileName {
	if v := os.Getenv("H3S_CREDENTIALS"); v != "" {
		return FileName(v)
	}
	return FileName("h3s-secrets.yaml")
}

const ( // #nosec G101 -- These are configuration file names, not credentials
	KubeConfigFileName FileName = "h3s-kubeconfig.yaml" // KubeConfigFileName is the file name for the h3s cluster configuration file
)
