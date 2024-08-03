package kubeconfig

import (
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
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
