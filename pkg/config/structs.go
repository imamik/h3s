package config

type NetworkZone string

const (
	EUCentral NetworkZone = "eu-central"
	USEast    NetworkZone = "us-east"
	USWest    NetworkZone = "us-west"
)

type Location string

const (
	Nuremberg   Location = "nbg1" // DE - EU Central
	Falkenstein Location = "fsn1" // DE - EU Central
	Helsinki    Location = "hel1" // FI - EU Central
	Ashburn     Location = "ash"  // US - US East
	Hillsboro   Location = "hil"  // US - US West
)

type NodePool struct {
	Name     string   `yaml:"name"`
	nodes    int      `yaml:"nodes"`
	Location Location `yaml:"location"`
}

type ControlPlanePool struct {
	NodePool
	AsWorkerpool bool `yaml:"asWorker"`
	LoadBalancer bool `yaml:"loadBalancer"`
}

type WorkerPools struct {
	LoadBalancer bool       `yaml:"loadBalancer"`
	WorkerPools  []NodePool `yaml:"workerPools"`
}

type Config struct {
	Name                 string           `yaml:"name"`
	K3sVersion           string           `yaml:"k3sVersion"`
	NetworkZone          NetworkZone      `yaml:"region"`
	ControlPlanePool     ControlPlanePool `yaml:"controlPlanePool"`
	WorkerPools          WorkerPools      `yaml:"workerPools"`
	CombinedLoadBalancer bool             `yaml:"combinedLoadBalancer"`
}
