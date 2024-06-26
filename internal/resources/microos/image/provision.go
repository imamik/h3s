package image

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos/image/commands"
)

func Provision(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
	server *hcloud.Server,
) {
	execute(ctx, server, commands.DownloadImage(architecture), 0, false)
	execute(ctx, server, commands.WriteImage(), 0, true)
	execute(ctx, server, commands.Packages(), 5, true)
	execute(ctx, server, commands.CleanUp(), 5, false)
}

func execute(ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	cmd string,
	pauseBeforeSeconds int,
	expectDisconnect bool,
) {
	fmt.Printf(`
=============== Executing command on server ===============
%s
===========================================================
`, cmd)
}
