package kubeconfig

// KubeConfig represents the overall Kubernetes configuration
type KubeConfig struct {
	APIVersion     string      `yaml:"apiVersion"`
	Kind           string      `yaml:"kind"`
	CurrentContext string      `yaml:"current-context"`
	Clusters       []Cluster   `yaml:"clusters"`
	Contexts       []Context   `yaml:"contexts"`
	Users          []User      `yaml:"users"`
	Preferences    Preferences `yaml:"preferences,omitempty"`
}

// Preferences represents user-specific preferences
type Preferences struct {
	Colors bool `yaml:"colors,omitempty"`
}

// Cluster represents a Kubernetes cluster configuration
type Cluster struct {
	Name    string         `yaml:"name"`
	Cluster ClusterDetails `yaml:"cluster"`
}

// ClusterDetails contains the details of the cluster
type ClusterDetails struct {
	Server                   string `yaml:"server"`
	CertificateAuthority     string `yaml:"certificate-authority,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

// Context represents a Kubernetes context configuration
type Context struct {
	Name    string         `yaml:"name"`
	Context ContextDetails `yaml:"context"`
}

// ContextDetails contains the details of the context
type ContextDetails struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

// User represents a Kubernetes user configuration
type User struct {
	Name string      `yaml:"name"`
	User UserDetails `yaml:"user"`
}

// UserDetails contains the details of the user credentials
type UserDetails struct {
	ClientCertificate     string `yaml:"client-certificate,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKey             string `yaml:"client-key,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Token                 string `yaml:"token,omitempty"`
	Username              string `yaml:"username,omitempty"`
	Password              string `yaml:"password,omitempty"`
}
