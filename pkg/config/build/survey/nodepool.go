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

	huh.NewConfirm().
		Title("Enable IPv4").
		Description("Enable IPv4 for the nodes in this pool").
		Value(&nodePool.EnableIPv4).
		Run()

	huh.NewConfirm().
		Title("Enable IPv6").
		Description("Enable IPv6 for the nodes in this pool").
		Value(&nodePool.EnableIPv6).
		Run()

	nodePool.Location = getLocation(
		"Location",
		fmt.Sprintf("Number of nodes in the '%s' pool ", nodePool.Name),
		networkZone,
	)

	nodePool.Instance = getInstance()

	return nodePool
}
