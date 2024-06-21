package survey

import (
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/k3s/releases"
	"strconv"
)

func Survey(k3sReleases []releases.Release) (config.Config, error) {
	var conf config.Config

	huh.NewInput().
		Title("Project Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&conf.Name).
		Run()

	huh.NewSelect[string]().
		Title("K3S Version").
		Description("The version of K3s to install").
		Options(releasesToOptions(k3sReleases)...).
		Value(&conf.K3sVersion).
		Run()

	huh.NewSelect[hcloud.NetworkZone]().
		Title("Network Zone").
		Description("The network zone to deploy the cluster in").
		Options(networkZoneOptions...).
		Value(&conf.NetworkZone).
		Run()

	conf.ControlPlane.Location = getLocation("Control Plane Location", "Location of the control plane node", conf.NetworkZone)

	var controlPlaneNodesString string
	huh.NewInput().
		Title("Nodes").
		Description("Number of control plane controlPlaneNodes. Must be an uneven number").
		Value(&controlPlaneNodesString).
		Validate(config.IsUnevenNumberString).
		CharLimit(1).
		Run()

	controlPlaneNodes, _ := strconv.Atoi(controlPlaneNodesString)
	conf.ControlPlane.Nodes = controlPlaneNodes

	huh.NewConfirm().
		Title("Control Plane Load Balancer").
		Description("Use a load balancer for the control plane controlPlaneNodes").
		Value(&conf.ControlPlane.LoadBalancer).
		Run()

	huh.NewConfirm().
		Title("Control Plane as Worker Pool").
		Description("Use the control plane controlPlaneNodes as workers").
		Value(&conf.ControlPlane.AsWorkerPool).
		Run()

	huh.NewConfirm().
		Title("Combined Load Balancer").
		Description("Use a single load balancer for all controlPlaneNodes").
		Value(&conf.CombinedLoadBalancer).
		Run()

	if conf.CombinedLoadBalancer {
		conf.ControlPlane.LoadBalancer = false
	}

	var workerPoolsString string
	huh.NewInput().
		Title("Worker Pools").
		Description("Number of worker pools").
		Value(&workerPoolsString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run()

	workerPools, _ := strconv.Atoi(workerPoolsString)
	conf.ControlPlane.Nodes = controlPlaneNodes

	for i := 0; i < workerPools; i++ {
		nodePool := getNodePool(conf.NetworkZone)
		conf.WorkerPools = append(conf.WorkerPools, nodePool)
	}

	return conf, nil

}
