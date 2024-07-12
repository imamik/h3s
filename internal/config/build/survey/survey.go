package survey

import (
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/k3s/releases"
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

	conf.SSHKeyPaths.PrivateKeyPath = "$HOME/.ssh/id_ed25519"

	huh.NewInput().
		Title("SHH Private Key Path").
		Description("Path to the private key to use for SSH").
		Value(&conf.SSHKeyPaths.PrivateKeyPath).
		Run()

	conf.SSHKeyPaths.PublicKeyPath = conf.SSHKeyPaths.PrivateKeyPath + ".pub"

	huh.NewInput().
		Title("Domain").
		Description("The domain you want to setup (e.g. 'example.com')").
		Value(&conf.Domain).
		Run()

	huh.NewSelect[hcloud.NetworkZone]().
		Title("Network Zone").
		Description("The network zone to deploy the cluster in").
		Options(networkZoneOptions...).
		Value(&conf.NetworkZone).
		Run()

	conf.ControlPlane.Pool.Name = "control-plane"
	conf.ControlPlane.Pool.Location = getLocation("Control Plane Location", "Location of the control plane node", conf.NetworkZone)
	conf.ControlPlane.Pool.Instance = getInstance()

	var controlPlaneNodesString string
	huh.NewInput().
		Title("Nodes").
		Description("Number of control plane nodes. Must be an uneven number").
		Value(&controlPlaneNodesString).
		Validate(config.IsUnevenNumberString).
		CharLimit(1).
		Run()

	controlPlaneNodes, _ := strconv.Atoi(controlPlaneNodesString)
	conf.ControlPlane.Pool.Nodes = controlPlaneNodes

	huh.NewConfirm().
		Title("Public IPs").
		Description("Should the nodes have public IPv4 & IPv6 (alt: all traffic is in private network & gateway is created)").
		Value(&conf.PublicIps).
		Run()

	huh.NewInput().
		Title("Certmanager Email").
		Description("Email to use for cert-manager (letsencrypt certificate generation process)").
		Value(&conf.CertManager.Email).
		Run()

	huh.NewConfirm().
		Title("Control Plane as Worker Pool").
		Description("Use the control plane controlPlaneNodes as workers").
		Value(&conf.ControlPlane.AsWorkerPool).
		Run()

	var workerPoolsString string
	huh.NewInput().
		Title("Worker Pools").
		Description("Number of worker pools").
		Value(&workerPoolsString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run()

	workerPools, _ := strconv.Atoi(workerPoolsString)
	conf.ControlPlane.Pool.Nodes = controlPlaneNodes

	for i := 0; i < workerPools; i++ {
		nodePool := getNodePool(conf.NetworkZone)
		conf.WorkerPools = append(conf.WorkerPools, nodePool)
	}

	return conf, nil

}
