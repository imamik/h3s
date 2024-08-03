package kubeconfig

import (
	"fmt"
	"h3s/internal/config/path"
	"h3s/internal/utils/yaml"
)

func SaveKubeConfig(projectName string, kubeConfig KubeConfig) {
	p := path.GetPath(projectName, path.KubeConfigFileName)
	err := yaml.Save(kubeConfig, p)
	if err != nil {
		fmt.Println(err)
		return
	}
}
