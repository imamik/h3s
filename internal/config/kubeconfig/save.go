package kubeconfig

import (
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
)

func SaveKubeConfig(projectName string, kubeConfig KubeConfig) error {
	p := string(path.KubeConfigFileName)
	_, err := file.New(p).SetYaml(kubeConfig).Save()
	return err
}
