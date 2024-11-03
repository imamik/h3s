package node

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

// Delete deletes a Hetzner cloud server
func Delete(
	ctx *cluster.Cluster,
	pool config.NodePool,
	i int,
) error {
	l := logger.New(nil, logger.Server, logger.Delete, getName(ctx, pool, i))
	defer l.LogEvents()

	server, err := Get(ctx, pool, i)

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
