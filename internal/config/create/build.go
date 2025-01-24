// Package create contains the functionality for creating a h3s cluster configuration
package create

import (
	"h3s/internal/config/create/survey"
	"h3s/internal/k3s"
	"h3s/internal/utils/file"
)

// Build surveys the user for the cluster configuration and saves it to h3s.yaml in the current directory.
func Build(k3sReleases []k3s.Release) error {
	conf, err := survey.Survey(k3sReleases)
	if err != nil {
		return err
	}

	if _, err = file.New("h3s.yaml").SetYaml(conf).Save(); err != nil {
		return err
	}

	return nil
}
