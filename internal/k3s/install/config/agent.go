package config

// K3sAgentConfig represents the configuration options for k3s agent command
type K3sAgentConfig struct {
	// Config and Logging
	ConfigFile      string `yaml:"config,omitempty"`          // Load configuration from FILE (default: "/etc/rancher/k3s/config.yaml")
	Debug           bool   `yaml:"debug,omitempty"`           // Turn on debug logs
	LogLevel        int    `yaml:"v,omitempty"`               // Number for the log level verbosity (default: 0)
	VModule         string `yaml:"vmodule,omitempty"`         // Comma-separated list of FILE_PATTERN=LOG_LEVEL settings for file-filtered logging
	LogFile         string `yaml:"log,omitempty"`             // Log to file
	AlsoLogToStderr bool   `yaml:"alsologtostderr,omitempty"` // Log to standard error as well as file (if set)

	// Cluster Join
	Token     string `yaml:"token,omitempty"`      // Token to use for authentication
	TokenFile string `yaml:"token-file,omitempty"` // Token file to use for authentication
	Server    string `yaml:"server,omitempty"`     // Server to connect to

	// data Directory
	DataDir string `yaml:"data-dir,omitempty"` // Folder to hold state (default: "/var/lib/rancher/k3s")

	// Node Configuration
	NodeName                      string   `yaml:"node-name,omitempty"`                         // Node name
	WithNodeID                    bool     `yaml:"with-node-id,omitempty"`                      // Append id to node name
	NodeLabel                     []string `yaml:"node-label,omitempty"`                        // Registering and starting kubelet with set of labels
	NodeTaint                     []string `yaml:"node-taint,omitempty"`                        // Registering kubelet with set of taints
	ImageCredentialProviderBinDir string   `yaml:"image-credential-provider-bin-dir,omitempty"` // The path to the directory where credential provider plugin binaries are located (default: "/var/lib/rancher/credentialprovider/bin")
	ImageCredentialProviderConfig string   `yaml:"image-credential-provider-config,omitempty"`  // The path to the credential provider plugin config file (default: "/var/lib/rancher/credentialprovider/config.yaml")
	SELinux                       bool     `yaml:"selinux,omitempty"`                           // Enable SELinux in containerd
	LBServerPort                  int      `yaml:"lb-server-port,omitempty"`                    // Local port for supervisor client load-balancer (default: 6444)
	ProtectKernelDefaults         bool     `yaml:"protect-kernel-defaults,omitempty"`           // Kernel tuning behavior. If set, error if kernel tunables are different than kubelet defaults.

	// Runtime Configuration
	ContainerRuntimeEndpoint string `yaml:"container-runtime-endpoint,omitempty"` // Disable embedded containerd and use the CRI socket at the given path; when used with --docker this sets the docker socket path
	PauseImage               string `yaml:"pause-image,omitempty"`                // Customized pause image for containerd or docker sandbox (default: "rancher/mirrored-pause:3.6")
	Snapshotter              string `yaml:"snapshotter,omitempty"`                // Override default containerd snapshotter (default: "overlayfs")
	PrivateRegistry          string `yaml:"private-registry,omitempty"`           // Private registry configuration file (default: "/etc/rancher/k3s/registries.yaml")
	Docker                   bool   `yaml:"docker,omitempty"`                     // (experimental) Use cri-dockerd instead of containerd

	// Agent Networking
	NodeIP         []string `yaml:"node-ip,omitempty"`          // IPv4/IPv6 addresses to advertise for node
	NodeExternalIP []string `yaml:"node-external-ip,omitempty"` // IPv4/IPv6 external IP addresses to advertise for node
	ResolvConf     string   `yaml:"resolv-conf,omitempty"`      // Kubelet resolv.conf file
	FlannelIface   string   `yaml:"flannel-iface,omitempty"`    // Override default flannel interface
	FlannelConf    string   `yaml:"flannel-conf,omitempty"`     // Override default flannel config file
	FlannelCNIConf string   `yaml:"flannel-cni-conf,omitempty"` // Override default flannel cni config file

	// Agent Flags
	KubeletArg   []string `yaml:"kubelet-arg,omitempty"`    // Customized flag for kubelet process
	KubeProxyArg []string `yaml:"kube-proxy-arg,omitempty"` // Customized flag for kube-proxy process

	// Experimental Features
	Rootless         bool `yaml:"rootless,omitempty"`           // Run rootless
	PreferBundledBin bool `yaml:"prefer-bundled-bin,omitempty"` // Prefer bundled userspace binaries over host binaries
}
