// Package server contains the functionality for managing Hetzner cloud servers
package server

import (
	"errors"
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/utils/logger"
	"sort"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// AllServers is a collection of all servers in the cluster, sorted by role (ControlPlane, Worker, Gateway, Other)
type AllServers struct {
	ControlPlane []*hcloud.Server
	Worker       []*hcloud.Server
	Gateway      []*hcloud.Server
	Other        []*hcloud.Server
}

func getAll(ctx *cluster.Cluster) ([]*hcloud.Server, error) {
	l := logger.New(nil, logger.Server, logger.Get, "All")
	defer l.LogEvents()

	label := fmt.Sprintf("project=%s", ctx.Config.Name)

	servers, err := ctx.CloudClient.Server.AllWithOpts(ctx.Context, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: label,
		},
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if len(servers) == 0 {
		err = errors.New("no servers found")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return servers, nil
}

// GetAll returns all servers in the cluster, sorted by name & grouped by role (ControlPlane, Worker, Gateway, Other)
func GetAll(ctx *cluster.Cluster) (*AllServers, error) {
	l := logger.New(nil, logger.Server, logger.Get, "All")
	defer l.LogEvents()

	var err error
	var nodes []*hcloud.Server
	for {
		nodes, err = getAll(ctx)
		if err != nil {
			l.AddEvent(logger.Failure, err)
			return nil, err
		}
		if len(nodes) == 0 {
			err = errors.New("no nodes found")
			l.AddEvent(logger.Failure, err)
			return nil, err
		}

		allNodesHavePrivateIP := true
		for _, n := range nodes {
			if len(n.PrivateNet) < 1 {
				allNodesHavePrivateIP = false
				break
			}
		}

		if allNodesHavePrivateIP {
			break
		}

		l.AddEvent(logger.Info, "Not all nodes have an assigned private IP. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	all := AllServers{}

	for _, n := range nodes {
		if node.IsControlPlaneNode(n) {
			all.ControlPlane = append(all.ControlPlane, n)
		} else if node.IsWorkerNode(n) {
			all.Worker = append(all.Worker, n)
		} else if node.IsGatewayNode(n) {
			all.Gateway = append(all.Gateway, n)
		} else {
			all.Other = append(all.Other, n)
		}
	}

	l.AddEvent(logger.Success)
	return &all, nil
}
