package credentials

type ProjectCredentials struct {
	HCloudToken string `yaml:"hcloud_token"`
	K3sToken    string `yaml:"k3s_token"`
}

type Credentials map[string]ProjectCredentials