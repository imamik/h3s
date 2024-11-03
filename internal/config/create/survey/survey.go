// Package survey contains the survey for creating a h3s cluster
package survey

import (
	"h3s/internal/config"
	"h3s/internal/k3s"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Survey surveys the user for the cluster configuration and returns it.
func Survey(k3sReleases []k3s.Release) (config.Config, error) {
	var conf config.Config

	if inputErr := huh.NewInput().
		Title("Project Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&conf.Name).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	if inputErr := huh.NewSelect[string]().
		Title("K3S Version").
		Description("The version of K3s to install").
		Options(releasesToOptions(k3sReleases)...).
		Value(&conf.K3sVersion).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	conf.SSHKeyPaths.PrivateKeyPath = "$HOME/.ssh/id_ed25519"

	if inputErr := huh.NewInput().
		Title("SHH Private Key Path").
		Description("Path to the private key to use for SSH").
		Value(&conf.SSHKeyPaths.PrivateKeyPath).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	conf.SSHKeyPaths.PublicKeyPath = conf.SSHKeyPaths.PrivateKeyPath + ".pub"

	if inputErr := huh.NewInput().
		Title("Domain").
		Description("The domain you want to setup (e.g. 'example.com')").
		Value(&conf.Domain).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	if inputErr := huh.NewSelect[hcloud.NetworkZone]().
		Title("Network Zone").
		Description("The network zone to deploy the cluster in").
		Options(networkZoneOptions...).
		Value(&conf.NetworkZone).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	conf.ControlPlane.Pool.Name = "control-plane"
	conf.ControlPlane.Pool.Location = getLocation("Control Plane Location", "Location of the control plane node", conf.NetworkZone)
	conf.ControlPlane.Pool.Instance = getInstance()

	var controlPlaneNodesString string
	if inputErr := huh.NewInput().
		Title("Nodes").
		Description("Number of control plane nodes. Must be an uneven number").
		Value(&controlPlaneNodesString).
		Validate(config.IsUnevenNumberString).
		CharLimit(1).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	controlPlaneNodes, err := strconv.Atoi(controlPlaneNodesString)
	if err != nil {
		return conf, err
	}
	conf.ControlPlane.Pool.Nodes = controlPlaneNodes

	if inputErr := huh.NewInput().
		Title("Certmanager Email").
		Description("Email to use for cert-manager (letsencrypt certificate generation process)").
		Value(&conf.CertManager.Email).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	if inputErr := huh.NewConfirm().
		Title("Control Plane as Worker Pool").
		Description("Use the control plane controlPlaneNodes as workers").
		Value(&conf.ControlPlane.AsWorkerPool).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	var workerPoolsString string
	if inputErr := huh.NewInput().
		Title("Worker Pools").
		Description("Number of worker pools").
		Value(&workerPoolsString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run(); inputErr != nil {
		return conf, inputErr
	}

	workerPools, err := strconv.Atoi(workerPoolsString)
	if err != nil {
		return conf, err
	}

	conf.ControlPlane.Pool.Nodes = controlPlaneNodes
	for i := 0; i < workerPools; i++ {
		nodePool, err := getNodePool(conf.NetworkZone)
		if err != nil {
			return conf, err
		}
		conf.WorkerPools = append(conf.WorkerPools, nodePool)
	}

	return conf, nil
}
