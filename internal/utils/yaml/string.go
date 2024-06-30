package yaml

import "gopkg.in/yaml.v3"

func String(config interface{}) string {
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		// handle error
		return ""
	}
	return string(yamlData)
}
