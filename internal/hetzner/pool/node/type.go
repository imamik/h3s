package node

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

type Type string

const (
	IsControlPlane Type = "is_control_plane"
	IsWorker       Type = "is_worker"
	IsGateway      Type = "is_gateway"
)

func IsControlPlaneNode(n *hcloud.Server) bool {
	return n.Labels[string(IsControlPlane)] == "true"
}

func IsWorkerNode(n *hcloud.Server) bool {
	return n.Labels[string(IsWorker)] == "true"
}

func IsGatewayNode(n *hcloud.Server) bool {
	return n.Labels[string(IsGateway)] == "true"
}
