package config

// CommonConfig represents the configuration options for k3s agents and servers
type CommonConfig struct {
	ResolvConf            string   `yaml:"resolv-conf,omitempty"`
	Server                string   `yaml:"server,omitempty"`
	FlannelCNIConf        string   `yaml:"flannel-cni-conf,omitempty"`
	FlannelConf           string   `yaml:"flannel-conf,omitempty"`
	FlannelIface          string   `yaml:"flannel-iface,omitempty"`
	Token                 string   `yaml:"token,omitempty"`
	VModule               string   `yaml:"vmodule,omitempty"`
	LogFile               string   `yaml:"log,omitempty"`
	ConfigFile            string   `yaml:"config,omitempty"`
	DataDir               string   `yaml:"data-dir,omitempty"`
	NodeName              string   `yaml:"node-name,omitempty"`
	TokenFile             string   `yaml:"token-file,omitempty"`
	NodeExternalIP        []string `yaml:"node-external-ip,omitempty"`
	NodeLabel             []string `yaml:"node-label,omitempty"`
	NodeIP                []string `yaml:"node-ip,omitempty"`
	KubeletArg            []string `yaml:"kubelet-arg,omitempty"`
	NodeTaint             []string `yaml:"node-taint,omitempty"`
	LBServerPort          int      `yaml:"lb-server-port,omitempty"`
	LogLevel              int      `yaml:"v,omitempty"`
	ProtectKernelDefaults bool     `yaml:"protect-kernel-defaults,omitempty"`
	SELinux               bool     `yaml:"selinux,omitempty"`
	Debug                 bool     `yaml:"debug,omitempty"`
	PreferBundledBin      bool     `yaml:"prefer-bundled-bin,omitempty"`
	AlsoLogToStderr       bool     `yaml:"alsologtostderr,omitempty"`
	Rootless              bool     `yaml:"rootless,omitempty"`
	WithNodeID            bool     `yaml:"with-node-id,omitempty"`
}

// AgentConfig represents the configuration options for k3s agent command
type AgentConfig struct {
	ImageCredentialProviderBinDir string `yaml:"image-credential-provider-bin-dir,omitempty"`
	ImageCredentialProviderConfig string `yaml:"image-credential-provider-config,omitempty"`
	ContainerRuntimeEndpoint      string `yaml:"container-runtime-endpoint,omitempty"`
	PauseImage                    string `yaml:"pause-image,omitempty"`
	Snapshotter                   string `yaml:"snapshotter,omitempty"`
	PrivateRegistry               string `yaml:"private-registry,omitempty"`
	CommonConfig                  `yaml:",inline"`
	Docker                        bool `yaml:"docker,omitempty"`
}

// ServerConfig represents the configuration options for k3s server command
type ServerConfig struct {
	DatastoreCAFile               string   `yaml:"datastore-cafile,omitempty"`
	EtcdS3SecretKey               string   `yaml:"etcd-s3-secret-key,omitempty"`
	DefaultLocalStoragePath       string   `yaml:"default-local-storage-path,omitempty"`
	AdvertiseAddress              string   `yaml:"advertise-address,omitempty"`
	EtcdS3Timeout                 string   `yaml:"etcd-s3-timeout,omitempty"`
	EtcdS3Folder                  string   `yaml:"etcd-s3-folder,omitempty"`
	EtcdS3Region                  string   `yaml:"etcd-s3-region,omitempty"`
	ServiceCIDR                   string   `yaml:"service-cidr,omitempty"`
	ServiceNodePortRange          string   `yaml:"service-node-port-range,omitempty"`
	ClusterDNS                    string   `yaml:"cluster-dns,omitempty"`
	ClusterDomain                 string   `yaml:"cluster-domain,omitempty"`
	FlannelBackend                string   `yaml:"flannel-backend,omitempty"`
	EtcdS3Bucket                  string   `yaml:"etcd-s3-bucket,omitempty"`
	BindAddress                   string   `yaml:"bind-address,omitempty"`
	EgressSelectorMode            string   `yaml:"egress-selector-mode,omitempty"`
	ServiceLBNamespace            string   `yaml:"servicelb-namespace,omitempty"`
	WriteKubeconfig               string   `yaml:"write-kubeconfig,omitempty"`
	WriteKubeconfigMode           string   `yaml:"write-kubeconfig-mode,omitempty"`
	AgentToken                    string   `yaml:"agent-token,omitempty"`
	AgentTokenFile                string   `yaml:"agent-token-file,omitempty"`
	EtcdS3AccessKey               string   `yaml:"etcd-s3-access-key,omitempty"`
	EtcdS3EndpointCA              string   `yaml:"etcd-s3-endpoint-ca,omitempty"`
	ClusterResetRestorePath       string   `yaml:"cluster-reset-restore-path,omitempty"`
	EtcdS3Endpoint                string   `yaml:"etcd-s3-endpoint,omitempty"`
	EtcdSnapshotDir               string   `yaml:"etcd-snapshot-dir,omitempty"`
	DatastoreCertFile             string   `yaml:"datastore-certfile,omitempty"`
	DatastoreEndpoint             string   `yaml:"datastore-endpoint,omitempty"`
	EtcdSnapshotScheduleCron      string   `yaml:"etcd-snapshot-schedule-cron,omitempty"`
	EtcdSnapshotName              string   `yaml:"etcd-snapshot-name,omitempty"`
	ClusterCIDR                   string   `yaml:"cluster-cidr,omitempty"`
	DatastoreKeyFile              string   `yaml:"datastore-keyfile,omitempty"`
	KubeSchedulerArg              []string `yaml:"kube-scheduler-arg,omitempty"`
	DisableComponents             []string `yaml:"disable,omitempty"`
	KubeCloudControllerManagerArg []string `yaml:"kube-cloud-controller-manager-arg,omitempty"`
	EtcdArg                       []string `yaml:"etcd-arg,omitempty"`
	KubeControllerManagerArg      []string `yaml:"kube-controller-manager-arg,omitempty"`
	TLSSAN                        []string `yaml:"tls-san,omitempty"`
	KubeAPIServerArg              []string `yaml:"kube-apiserver-arg,omitempty"`
	CommonConfig                  `yaml:",inline"`
	HTTPSListenPort               int  `yaml:"https-listen-port,omitempty"`
	AdvertisePort                 int  `yaml:"advertise-port,omitempty"`
	EtcdSnapshotRetention         int  `yaml:"etcd-snapshot-retention,omitempty"`
	EtcdSnapshotCompress          bool `yaml:"etcd-snapshot-compress,omitempty"`
	FlannelExternalIP             bool `yaml:"flannel-external-ip,omitempty"`
	EtcdExposeMetrics             bool `yaml:"etcd-expose-metrics,omitempty"`
	EtcdS3SkipSSLVerify           bool `yaml:"etcd-s3-skip-ssl-verify,omitempty"`
	ClusterReset                  bool `yaml:"cluster-reset,omitempty"`
	EtcdS3                        bool `yaml:"etcd-s3,omitempty"`
	EtcdS3Insecure                bool `yaml:"etcd-s3-insecure,omitempty"`
	FlannelIPv6Masq               bool `yaml:"flannel-ipv6-masq,omitempty"`
	EtcdDisableSnapshots          bool `yaml:"etcd-disable-snapshots,omitempty"`
	ClusterInit                   bool `yaml:"cluster-init,omitempty"`
	DisableScheduler              bool `yaml:"disable-scheduler,omitempty"`
	DisableCloudController        bool `yaml:"disable-cloud-controller,omitempty"`
	DisableKubeProxy              bool `yaml:"disable-kube-proxy,omitempty"`
	DisableNetworkPolicy          bool `yaml:"disable-network-policy,omitempty"`
	DisableHelmController         bool `yaml:"disable-helm-controller,omitempty"`
	SecretsEncryption             bool `yaml:"secrets-encryption,omitempty"`
	EnablePprof                   bool `yaml:"enable-pprof,omitempty"`
}
