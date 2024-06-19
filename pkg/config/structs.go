package config

type NodePool struct {
	Name     string   `yaml:"name"`
	Nodes    int      `yaml:"nodes"`
	Location Location `yaml:"location"`
}

type ControlPlanePool struct {
	Nodes        int      `yaml:"nodes"`
	Location     Location `yaml:"location"`
	AsWorkerPool bool     `yaml:"asWorkerPool"`
	LoadBalancer bool     `yaml:"loadBalancer"`
}

type Config struct {
	Name                 string           `yaml:"name"`
	K3sVersion           string           `yaml:"k3sVersion"`
	NetworkZone          NetworkZone      `yaml:"region"`
	ControlPlanePool     ControlPlanePool `yaml:"controlPlanePool"`
	WorkerPools          []NodePool       `yaml:"workerPools"`
	CombinedLoadBalancer bool             `yaml:"combinedLoadBalancer"`
}
