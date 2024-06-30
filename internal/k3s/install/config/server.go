package config

// K3sServerConfig represents the configuration options for k3s server command
type K3sServerConfig struct {
	// Config and Logging
	ConfigFile      string `yaml:"config,omitempty"`          // Load configuration from FILE (default: "/etc/rancher/k3s/config.yaml")
	Debug           bool   `yaml:"debug,omitempty"`           // Turn on debug logs
	LogLevel        int    `yaml:"v,omitempty"`               // Number for the log level verbosity (default: 0)
	VModule         string `yaml:"vmodule,omitempty"`         // Comma-separated list of FILE_PATTERN=LOG_LEVEL settings for file-filtered logging
	LogFile         string `yaml:"log,omitempty"`             // Log to file
	AlsoLogToStderr bool   `yaml:"alsologtostderr,omitempty"` // Log to standard error as well as file (if set)

	// Listener Configuration
	BindAddress      string   `yaml:"bind-address,omitempty"`      // k3s bind address (default: 0.0.0.0)
	HTTPSListenPort  int      `yaml:"https-listen-port,omitempty"` // HTTPS listen port (default: 6443)
	AdvertiseAddress string   `yaml:"advertise-address,omitempty"` // IPv4 address that apiserver uses to advertise to members of the cluster (default: node-external-ip/node-ip)
	AdvertisePort    int      `yaml:"advertise-port,omitempty"`    // Port that apiserver uses to advertise to members of the cluster (default: listen-port) (default: 0)
	TLSSAN           []string `yaml:"tls-san,omitempty"`           // Add additional hostnames or IPv4/IPv6 addresses as Subject Alternative Names on the server TLS cert

	// Data Directory
	DataDir string `yaml:"data-dir,omitempty"` // Folder to hold state (default: /var/lib/rancher/k3s or ${HOME}/.rancher/k3s if not root)

	// Networking Configuration
	ClusterCIDR          string `yaml:"cluster-cidr,omitempty"`            // IPv4/IPv6 network CIDRs to use for pod IPs (default: 10.42.0.0/16)
	ServiceCIDR          string `yaml:"service-cidr,omitempty"`            // IPv4/IPv6 network CIDRs to use for service IPs (default: 10.43.0.0/16)
	ServiceNodePortRange string `yaml:"service-node-port-range,omitempty"` // Port range to reserve for services with NodePort visibility (default: "30000-32767")
	ClusterDNS           string `yaml:"cluster-dns,omitempty"`             // IPv4 Cluster IP for coredns service. Should be in your service-cidr range (default: 10.43.0.10)
	ClusterDomain        string `yaml:"cluster-domain,omitempty"`          // Cluster Domain (default: "cluster.local")
	FlannelBackend       string `yaml:"flannel-backend,omitempty"`         // Backend configuration for Flannel networking (default: "vxlan")
	FlannelIPv6Masq      bool   `yaml:"flannel-ipv6-masq,omitempty"`       // Enable IPv6 masquerading for pod
	FlannelExternalIP    bool   `yaml:"flannel-external-ip,omitempty"`     // Use node external IP addresses for Flannel traffic
	EgressSelectorMode   string `yaml:"egress-selector-mode,omitempty"`    // One of 'agent', 'cluster', 'pod', 'disabled' (default: "agent")
	ServiceLBNamespace   string `yaml:"servicelb-namespace,omitempty"`     // Namespace of the pods for the servicelb component (default: "kube-system")

	// Client Configuration
	WriteKubeconfig     string `yaml:"write-kubeconfig,omitempty"`      // Write kubeconfig for admin client to this file
	WriteKubeconfigMode string `yaml:"write-kubeconfig-mode,omitempty"` // Write kubeconfig with this mode

	// Cluster Join
	Token          string `yaml:"token,omitempty"`            // Shared secret used to join a server or agent to a cluster
	TokenFile      string `yaml:"token-file,omitempty"`       // File containing the token
	AgentToken     string `yaml:"agent-token,omitempty"`      // Shared secret used to join agents to the cluster, but not servers
	AgentTokenFile string `yaml:"agent-token-file,omitempty"` // File containing the agent secret
	Server         string `yaml:"server,omitempty"`           // Server to connect to, used to join a cluster

	// Cluster Management
	ClusterInit             bool   `yaml:"cluster-init,omitempty"`               // Initialize a new cluster using embedded Etcd
	ClusterReset            bool   `yaml:"cluster-reset,omitempty"`              // Forget all peers and become sole member of a new cluster
	ClusterResetRestorePath string `yaml:"cluster-reset-restore-path,omitempty"` // Path to snapshot file to be restored

	// Process Customization
	KubeAPIServerArg              []string `yaml:"kube-apiserver-arg,omitempty"`                // Customized flag for kube-apiserver process
	EtcdArg                       []string `yaml:"etcd-arg,omitempty"`                          // Customized flag for etcd process
	KubeControllerManagerArg      []string `yaml:"kube-controller-manager-arg,omitempty"`       // Customized flag for kube-controller-manager process
	KubeSchedulerArg              []string `yaml:"kube-scheduler-arg,omitempty"`                // Customized flag for kube-scheduler process
	KubeCloudControllerManagerArg []string `yaml:"kube-cloud-controller-manager-arg,omitempty"` // Customized flag for kube-cloud-controller-manager process

	// Datastore Configuration
	DatastoreEndpoint        string `yaml:"datastore-endpoint,omitempty"`          // Specify etcd, Mysql, Postgres, or Sqlite (default) data source name
	DatastoreCAFile          string `yaml:"datastore-cafile,omitempty"`            // TLS Certificate Authority file used to secure datastore backend communication
	DatastoreCertFile        string `yaml:"datastore-certfile,omitempty"`          // TLS certification file used to secure datastore backend communication
	DatastoreKeyFile         string `yaml:"datastore-keyfile,omitempty"`           // TLS key file used to secure datastore backend communication
	EtcdExposeMetrics        bool   `yaml:"etcd-expose-metrics,omitempty"`         // Expose etcd metrics to client interface
	EtcdDisableSnapshots     bool   `yaml:"etcd-disable-snapshots,omitempty"`      // Disable automatic etcd snapshots
	EtcdSnapshotName         string `yaml:"etcd-snapshot-name,omitempty"`          // Set the base name of etcd snapshots (default: etcd-snapshot-<unix-timestamp>)
	EtcdSnapshotScheduleCron string `yaml:"etcd-snapshot-schedule-cron,omitempty"` // Snapshot interval time in cron spec. eg. every 5 hours '* */5 * * *'
	EtcdSnapshotRetention    int    `yaml:"etcd-snapshot-retention,omitempty"`     // Number of snapshots to retain (default: 5)
	EtcdSnapshotDir          string `yaml:"etcd-snapshot-dir,omitempty"`           // Directory to save db snapshots. (default: ${data-dir}/db/snapshots)
	EtcdSnapshotCompress     bool   `yaml:"etcd-snapshot-compress,omitempty"`      // Compress etcd snapshot

	// S3 Backup Configuration
	EtcdS3              bool   `yaml:"etcd-s3,omitempty"`                 // Enable backup to S3
	EtcdS3Endpoint      string `yaml:"etcd-s3-endpoint,omitempty"`        // S3 endpoint url (default: "s3.amazonaws.com")
	EtcdS3EndpointCA    string `yaml:"etcd-s3-endpoint-ca,omitempty"`     // S3 custom CA cert to connect to S3 endpoint
	EtcdS3SkipSSLVerify bool   `yaml:"etcd-s3-skip-ssl-verify,omitempty"` // Disables S3 SSL certificate validation
	EtcdS3AccessKey     string `yaml:"etcd-s3-access-key,omitempty"`      // S3 access key
	EtcdS3SecretKey     string `yaml:"etcd-s3-secret-key,omitempty"`      // S3 secret key
	EtcdS3Bucket        string `yaml:"etcd-s3-bucket,omitempty"`          // S3 bucket name
	EtcdS3Region        string `yaml:"etcd-s3-region,omitempty"`          // S3 region / bucket location (optional) (default: "us-east-1")
	EtcdS3Folder        string `yaml:"etcd-s3-folder,omitempty"`          // S3 folder
	EtcdS3Insecure      bool   `yaml:"etcd-s3-insecure,omitempty"`        // Disables S3 over HTTPS
	EtcdS3Timeout       string `yaml:"etcd-s3-timeout,omitempty"`         // S3 timeout (default: 5m0s)

	// Local Storage Configuration
	DefaultLocalStoragePath string `yaml:"default-local-storage-path,omitempty"` // Default local storage path for local provisioner storage class

	// Component Management
	DisableComponents      []string `yaml:"disable,omitempty"`                  // Do not deploy packaged components and delete any deployed components
	DisableScheduler       bool     `yaml:"disable-scheduler,omitempty"`        // Disable Kubernetes default scheduler
	DisableCloudController bool     `yaml:"disable-cloud-controller,omitempty"` // Disable k3s default cloud controller manager
	DisableKubeProxy       bool     `yaml:"disable-kube-proxy,omitempty"`       // Disable running kube-proxy
	DisableNetworkPolicy   bool     `yaml:"disable-network-policy,omitempty"`   // Disable k3s default network policy controller
	DisableHelmController  bool     `yaml:"disable-helm-controller,omitempty"`  // Disable Helm controller

	// Node Configuration
	NodeName   string   `yaml:"node-name,omitempty"`    // Node name
	WithNodeID bool     `yaml:"with-node-id,omitempty"` // Append id to node name
	NodeLabel  []string `yaml:"node-label,omitempty"`   // Registering and starting kubelet with set of labels
	NodeTaint  []string `yaml:"node-taint,omitempty"`   // Registering kubelet with set of taints

	// Image Credential Provider
	ImageCredentialProviderBinDir string `yaml:"image-credential-provider-bin-dir,omitempty"` // The path to the directory where credential provider plugin binaries are located (default: "/var/lib/rancher/credentialprovider/bin")
	ImageCredentialProviderConfig string `yaml:"image-credential-provider-config,omitempty"`  // The path to the credential provider plugin config file (default: "/var/lib/rancher/credentialprovider/config.yaml")

	// Runtime Configuration
	Docker                   bool   `yaml:"docker,omitempty"`                     // (experimental) Use cri-dockerd instead of containerd
	ContainerRuntimeEndpoint string `yaml:"container-runtime-endpoint,omitempty"` // Disable embedded containerd and use the CRI socket at the given path; when used with --docker this sets the docker socket path
	PauseImage               string `yaml:"pause-image,omitempty"`                // Customized pause image for containerd or docker sandbox (default: "rancher/mirrored-pause:3.6")
	Snapshotter              string `yaml:"snapshotter,omitempty"`                // Override default containerd snapshotter (default: "overlayfs")
	PrivateRegistry          string `yaml:"private-registry,omitempty"`           // Private registry configuration file (default: "/etc/rancher/k3s/registries.yaml")
	SystemDefaultRegistry    string `yaml:"system-default-registry,omitempty"`    // Private registry to be used for all system images

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

	// Security
	ProtectKernelDefaults bool `yaml:"protect-kernel-defaults,omitempty"` // Kernel tuning behavior. If set, error if kernel tunables are different than kubelet defaults.
	SecretsEncryption     bool `yaml:"secrets-encryption,omitempty"`      // Enable secret encryption at rest
	SELinux               bool `yaml:"selinux,omitempty"`                 // Enable SELinux in containerd

	// Experimental Features
	EnablePprof      bool `yaml:"enable-pprof,omitempty"`       // Enable pprof endpoint on supervisor port
	Rootless         bool `yaml:"rootless,omitempty"`           // Run rootless
	PreferBundledBin bool `yaml:"prefer-bundled-bin,omitempty"` // Prefer bundled userspace binaries over host binaries

	// Load Balancer
	LBServerPort int `yaml:"lb-server-port,omitempty"` // Local port for supervisor client load-balancer (default: 6444)
}
