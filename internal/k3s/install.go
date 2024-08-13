package k3s

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/server"
	commands2 "h3s/internal/k3s/commands"
)

func Install(ctx *cluster.Cluster) error {
	lb := loadbalancers.Get(ctx)
	gatewayNode, _ := gateway.GetIfNeeded(ctx)
	nodes, err := server.GetAll(ctx)

	if err != nil {
		return err
	}

	for _, remote := range nodes.ControlPlane {
		commands2.ControlPlane(ctx, lb, nodes.ControlPlane, gatewayNode, remote)
	}

	for _, remote := range nodes.Worker {
		commands2.Worker(ctx, nodes.ControlPlane[0], gatewayNode, remote)
	}

	return nil
}
