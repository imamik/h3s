package credentials

import (
	"fmt"
	"h3s/internal/config"
)

func Get(conf config.Config) (*ProjectCredentials, error) {
	credentials, err := initialize(conf.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return credentials, nil
}
