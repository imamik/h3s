package kubeconfig

import (
	"fmt"
	"hcloud-k3s-cli/internal/config/path"
	"hcloud-k3s-cli/internal/utils/yaml"
)

func SaveKubeConfig(projectName string, kubeConfig KubeConfig) {
	p := path.GetPath(projectName, path.KubeConfigFileName)
	err := yaml.Save(kubeConfig, p)
	if err != nil {
		fmt.Println(err)
		return
	}
}
