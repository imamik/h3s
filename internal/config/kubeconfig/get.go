package kubeconfig

import (
	"hcloud-k3s-cli/internal/config/path"
	"hcloud-k3s-cli/internal/utils/file"
)

func GetPathIfExists(projectName string) (string, bool) {
	p := path.GetPath(projectName, path.KubeConfigFileName)
	if file.Exists(p) {
		normalizedPath, err := file.Normalize(p)
		if err != nil {
			return "", false
		}
		return normalizedPath, true
	}
	return "", false
}
