package kubeconfig

import (
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
)

func GetPathIfExists() (string, bool) {
	p := string(path.KubeConfigFileName)
	f := file.New(p)
	absPath, _ := f.Path()
	return absPath, f.Exists()
}
