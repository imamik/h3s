package gateway

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets the Hetzner cloud gateway
func Get(ctx *cluster.Cluster) (*hcloud.Server, error) {
	name := getName(ctx)

	l := logger.New(nil, logger.Server, logger.Get, name)
	defer l.LogEvents()

	// Get server by name
	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, name)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Check if server is nil
	if server == nil {
		err = nil
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	return server, nil
}
