package survey

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"strconv"
)

func getNodePool(networkZone hcloud.NetworkZone) config.NodePool {
	var nodePool config.NodePool

	huh.NewInput().
		Title("Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&nodePool.Name).
		Run()

	var nodePoolNodesString string
	huh.NewInput().
		Title("Nodes").
		Description(fmt.Sprintf("Number of nodes in the '%s' pool ", nodePool.Name)).
		Value(&nodePoolNodesString).
		Validate(config.IsNumberString).
		CharLimit(1).
		Run()

	nodePoolNodes, _ := strconv.Atoi(nodePoolNodesString)
	nodePool.Nodes = nodePoolNodes

	nodePool.Location = getLocation(
		"Location",
		fmt.Sprintf("Number of nodes in the '%s' pool ", nodePool.Name),
		networkZone,
	)

	nodePool.InstanceType = getInstance()

	return nodePool
}
