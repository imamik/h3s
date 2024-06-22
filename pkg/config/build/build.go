package build

import (
	"fmt"
	"hcloud-k3s-cli/pkg/config/build/survey"
	"hcloud-k3s-cli/pkg/k3s/releases"
	"hcloud-k3s-cli/pkg/utils/yaml"
)

func Build(k3sReleases []releases.Release) {

	conf, err := survey.Survey(k3sReleases)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Save(conf, "hcloud-k3s.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
}
