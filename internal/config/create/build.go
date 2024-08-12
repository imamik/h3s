package create

import (
	"fmt"
	"h3s/internal/config/create/survey"
	"h3s/internal/k3s"
	"h3s/internal/utils/file"
)

func Build(k3sReleases []k3s.Release) {

	conf, err := survey.Survey(k3sReleases)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = file.New("h3s.yaml").SetYaml(conf).Save()
	if err != nil {
		fmt.Println(err)
		return
	}
}
