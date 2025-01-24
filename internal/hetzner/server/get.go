// Package server contains functionality for managing Hetzner cloud servers
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

// Client represents a Hetzner server client
type Client struct {
	cluster *cluster.Cluster
}

// NewClient creates a new Hetzner server client
func NewClient(cluster *cluster.Cluster) *Client {
	return &Client{
		cluster: cluster,
	}
}

// AllServers is a collection of all servers in the cluster, sorted by role (ControlPlane, Worker, Gateway, Other)
type AllServers struct {
	ControlPlane []*hcloud.Server
	Worker       []*hcloud.Server
	Gateway      []*hcloud.Server
	Other        []*hcloud.Server
}

// GetAll returns all servers in the cluster, sorted by name & grouped by role (ControlPlane, Worker, Gateway, Other)
func GetAll(ctx *cluster.Cluster) (*AllServers, error) {
	client := NewClient(ctx)
	return client.GetAll()
}

func (c *Client) getAll() ([]*hcloud.Server, error) {
	l := logger.New(nil, logger.Server, logger.Get, "All")
	defer l.LogEvents()

	label := fmt.Sprintf("project=%s", c.cluster.Config.Name)

	servers, err := c.cluster.CloudClient.Server.AllWithOpts(c.cluster.Context, hcloud.ServerListOpts{
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
func (c *Client) GetAll() (*AllServers, error) {
	l := logger.New(nil, logger.Server, logger.Get, "All")
	defer l.LogEvents()

	nodes, err := c.fetchNodesWithPrivateIP(l)
	if err != nil {
		return nil, err
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	all := c.groupNodesByRole(nodes)

	l.AddEvent(logger.Success)
	return &all, nil
}

func (c *Client) fetchNodesWithPrivateIP(l *logger.EventLogger) ([]*hcloud.Server, error) {
	for {
		nodes, err := c.getAll()
		if err != nil {
			l.AddEvent(logger.Failure, err)
			return nil, err
		}
		if len(nodes) == 0 {
			err = errors.New("no nodes found")
			l.AddEvent(logger.Failure, err)
			return nil, err
		}

		if c.allNodesHavePrivateIP(nodes) {
			return nodes, nil
		}

		l.AddEvent(logger.Info, "Not all nodes have an assigned private IP. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}

func (c *Client) allNodesHavePrivateIP(nodes []*hcloud.Server) bool {
	for _, n := range nodes {
		if len(n.PrivateNet) < 1 {
			return false
		}
	}
	return true
}

func (c *Client) groupNodesByRole(nodes []*hcloud.Server) AllServers {
	all := AllServers{}

	for _, n := range nodes {
		switch {
		case node.IsControlPlaneNode(n):
			all.ControlPlane = append(all.ControlPlane, n)
		case node.IsWorkerNode(n):
			all.Worker = append(all.Worker, n)
		case node.IsGatewayNode(n):
			all.Gateway = append(all.Gateway, n)
		default:
			all.Other = append(all.Other, n)
		}
	}

	return all
}
