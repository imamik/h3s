package config

type NodePool struct {
	Name     string   `yaml:"name"`
	Nodes    int      `yaml:"nodes"`
	Location Location `yaml:"location"`
}

type ControlPlane struct {
	Nodes        int      `yaml:"nodes"`
	Location     Location `yaml:"location"`
	AsWorkerPool bool     `yaml:"asWorkerPool,omitempty"`
	LoadBalancer bool     `yaml:"loadBalancer,omitempty"`
}

type Config struct {
	Name                 string       `yaml:"name"`
	K3sVersion           string       `yaml:"k3sVersion"`
	NetworkZone          NetworkZone  `yaml:"networkZone"`
	ControlPlane         ControlPlane `yaml:"controlPlane"`
	WorkerPools          []NodePool   `yaml:"workerPools"`
	CombinedLoadBalancer bool         `yaml:"combinedLoadBalancer"`
}
