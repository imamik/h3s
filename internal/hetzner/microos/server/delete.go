package server

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Delete deletes a Hetzner cloud microOS server
func Delete(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
) error {
	l := logger.New(nil, logger.Server, logger.Delete, getName(ctx, architecture))
	defer l.LogEvents()

	server, err := Get(ctx, architecture)

	if server == nil && err.Error() == "server is nil" {
		l.AddEvent(logger.Success, "no server found to delete")
		return nil
	}

	_, _, err = ctx.CloudClient.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
