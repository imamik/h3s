package node

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// Type is the type of a Hetzner cloud server - can be control plane, worker or gateway
type Type string

const (
	IsControlPlane Type = "is_control_plane" // IsControlPlane is the label for a control plane node
	IsWorker       Type = "is_worker"        // IsWorker is the label for a worker node
	IsGateway      Type = "is_gateway"       // IsGateway is the label for a gateway node
)

const trueString = "true"

// IsControlPlaneNode returns true if the server is a control plane node
func IsControlPlaneNode(n *hcloud.Server) bool {
	return n.Labels[string(IsControlPlane)] == trueString
}

// IsWorkerNode returns true if the server is a worker node
func IsWorkerNode(n *hcloud.Server) bool {
	return n.Labels[string(IsWorker)] == trueString
}

// IsGatewayNode returns true if the server is a gateway node
func IsGatewayNode(n *hcloud.Server) bool {
	return n.Labels[string(IsGateway)] == trueString
}
