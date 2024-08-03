package build

import (
	"fmt"
	"h3s/internal/config/build/survey"
	"h3s/internal/k3s/releases"
	"h3s/internal/utils/yaml"
)

func Build(k3sReleases []releases.Release) {

	conf, err := survey.Survey(k3sReleases)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Save(conf, "h3s.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
}
