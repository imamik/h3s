package survey

import (
	"fmt"
	"h3s/internal/config"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getNodePool(networkZone hcloud.NetworkZone) (config.NodePool, error) {
	var nodePool config.NodePool

	err := huh.NewInput().
		Title("Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&nodePool.Name).
		Run()
	if err != nil {
		return config.NodePool{}, err
	}

	var nodePoolNodesString string
	err = huh.NewInput().
		Title("Nodes").
		Description(fmt.Sprintf("Number of nodes in the '%s' pool ", nodePool.Name)).
		Value(&nodePoolNodesString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run()
	if err != nil {
		return config.NodePool{}, err
	}

	nodePoolNodes, err := strconv.Atoi(nodePoolNodesString)
	if err != nil {
		return config.NodePool{}, err
	}
	nodePool.Nodes = nodePoolNodes

	nodePool.Location, err = getLocation(
		"Location",
		fmt.Sprintf("Number of nodes in the '%s' pool ", nodePool.Name),
		networkZone,
	)
	if err != nil {
		return config.NodePool{}, err
	}

	nodePool.Instance = getInstance()

	return nodePool, nil
}
