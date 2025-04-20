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

	// Collect basic configuration
	if err := surveyBasicConfig(&conf, k3sReleases); err != nil {
		return conf, err
	}

	// Configure control plane
	if err := setupControlPlane(&conf); err != nil {
		return conf, err
	}

	// Configure cert manager
	if err := surveyCertManagerEmail(&conf); err != nil {
		return conf, err
	}

	// Configure worker pools
	if err := configureWorkerPools(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

// surveyBasicConfig collects the basic configuration for the cluster
func surveyBasicConfig(conf *config.Config, k3sReleases []k3s.Release) error {
	// Project name
	if err := surveyProjectName(conf); err != nil {
		return err
	}

	// K3s version
	if err := surveyK3sVersion(k3sReleases, conf); err != nil {
		return err
	}

	// SSH key paths
	conf.SSHKeyPaths.PrivateKeyPath = "$HOME/.ssh/id_ed25519"
	if err := surveySSHPrivateKeyPath(conf); err != nil {
		return err
	}
	conf.SSHKeyPaths.PublicKeyPath = conf.SSHKeyPaths.PrivateKeyPath + ".pub"

	// Domain and network zone
	if err := surveyDomain(conf); err != nil {
		return err
	}

	if err := surveyNetworkZone(conf); err != nil {
		return err
	}

	return nil
}

// setupControlPlane configures the control plane
func setupControlPlane(conf *config.Config) error {
	// Set up control plane configuration
	conf.ControlPlane.Pool.Name = "control-plane"
	location, err := getLocation("Control Plane Location", "Location of the control plane node", conf.NetworkZone)
	if err != nil {
		return err
	}
	conf.ControlPlane.Pool.Location = location
	conf.ControlPlane.Pool.Instance = getInstance()

	// Control plane nodes
	if err := surveyControlPlaneNodes(conf); err != nil {
		return err
	}

	// Control plane as worker pool
	if err := surveyControlPlaneAsWorkerPool(conf); err != nil {
		return err
	}

	return nil
}

// configureWorkerPools sets up the worker pools
func configureWorkerPools(conf *config.Config) error {
	return surveyWorkerPools(conf)
}

func surveyProjectName(conf *config.Config) error {
	return huh.NewInput().
		Title("Project Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&conf.Name).
		Run()
}

func surveyK3sVersion(k3sReleases []k3s.Release, conf *config.Config) error {
	return huh.NewSelect[string]().
		Title("K3S Version").
		Description("The version of K3s to install").
		Options(releasesToOptions(k3sReleases)...).
		Value(&conf.K3sVersion).
		Run()
}

func surveySSHPrivateKeyPath(conf *config.Config) error {
	return huh.NewInput().
		Title("SHH Private Key Path").
		Description("Path to the private key to use for SSH").
		Value(&conf.SSHKeyPaths.PrivateKeyPath).
		Run()
}

func surveyDomain(conf *config.Config) error {
	return huh.NewInput().
		Title("Domain").
		Description("The domain you want to setup (e.g. 'example.com')").
		Value(&conf.Domain).
		Run()
}

func surveyNetworkZone(conf *config.Config) error {
	return huh.NewSelect[hcloud.NetworkZone]().
		Title("Network Zone").
		Description("The network zone to deploy the cluster in").
		Options(networkZoneOptions...).
		Value(&conf.NetworkZone).
		Run()
}

func surveyControlPlaneNodes(conf *config.Config) error {
	var controlPlaneNodesString string
	if err := huh.NewInput().
		Title("Nodes").
		Description("Number of control plane nodes. Must be an uneven number").
		Value(&controlPlaneNodesString).
		Validate(config.IsUnevenNumberString).
		CharLimit(1).
		Run(); err != nil {
		return err
	}

	controlPlaneNodes, err := strconv.Atoi(controlPlaneNodesString)
	if err != nil {
		return err
	}
	conf.ControlPlane.Pool.Nodes = controlPlaneNodes
	return nil
}

func surveyCertManagerEmail(conf *config.Config) error {
	return huh.NewInput().
		Title("Certmanager Email").
		Description("Email to use for cert-manager (letsencrypt certificate generation process)").
		Value(&conf.CertManager.Email).
		Run()
}

func surveyControlPlaneAsWorkerPool(conf *config.Config) error {
	return huh.NewConfirm().
		Title("Control Plane as Worker Pool").
		Description("Use the control plane controlPlaneNodes as workers").
		Value(&conf.ControlPlane.AsWorkerPool).
		Run()
}

func surveyWorkerPools(conf *config.Config) error {
	var workerPoolsString string
	if err := huh.NewInput().
		Title("Worker Pools").
		Description("Number of worker pools").
		Value(&workerPoolsString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run(); err != nil {
		return err
	}

	workerPools, err := strconv.Atoi(workerPoolsString)
	if err != nil {
		return err
	}

	for i := 0; i < workerPools; i++ {
		nodePool, err := getNodePool(conf.NetworkZone)
		if err != nil {
			return err
		}
		conf.WorkerPools = append(conf.WorkerPools, nodePool)
	}
	return nil
}
