package credentials

import (
	"fmt"
	"hcloud-k3s-cli/internal/config"
	"os"
)

func Get(conf config.Config) (ProjectCredentials, error) {
	credentials, err := initialize()
	if err != nil {
		fmt.Println(err)
		return ProjectCredentials{}, err
	}
	projectCredentials, ok := credentials[conf.Name]
	if ok {
		return projectCredentials, nil
	}

	fallbackCredentials := ProjectCredentials{
		HCloudToken: os.Getenv("HCLOUD_TOKEN"),
	}
	return fallbackCredentials, fmt.Errorf("project not found in credentials file")
}
