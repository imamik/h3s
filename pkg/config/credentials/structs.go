package credentials

type ProjectCredentials struct {
	HCloudToken string `yaml:"hcloud_token"`
}

type Credentials map[string]ProjectCredentials
