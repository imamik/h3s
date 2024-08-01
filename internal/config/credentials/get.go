package credentials

import (
	"fmt"
	"hcloud-k3s-cli/internal/config"
)

func Get(conf config.Config) (*ProjectCredentials, error) {
	credentials, err := initialize(conf.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return credentials, nil
}
