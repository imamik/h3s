package server

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Action represents a server operation type
type Action string

const (
	ActionShutdown   Action = "ShutDown"        // ActionShutdown will shut down the server
	ActionReboot     Action = "Reboot"          // ActionReboot will reboot the server
	ActionRescueMode Action = "Set Rescue Mode" // ActionRescueMode will set the rescue mode for the server
)

// performServerAction executes a server action and handles common patterns
func performServerAction(
	ctx *cluster.Cluster,
	server *hcloud.Server,
	actionType Action,
	actionFn func() (*hcloud.Action, *hcloud.Response, error),
) error {
	l := logger.New(nil, logger.Server, string(actionType), server.Name)
	defer l.LogEvents()

	action, _, err := actionFn()
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, action); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}

// Shutdown shuts down the given server.
func Shutdown(ctx *cluster.Cluster, server *hcloud.Server) error {
	return performServerAction(ctx, server, ActionShutdown, func() (*hcloud.Action, *hcloud.Response, error) {
		return ctx.CloudClient.Server.Shutdown(ctx.Context, server)
	})
}

// Reboot reboots the Hetzner cloud microOS server
func Reboot(ctx *cluster.Cluster, server *hcloud.Server) error {
	return performServerAction(ctx, server, ActionReboot, func() (*hcloud.Action, *hcloud.Response, error) {
		return ctx.CloudClient.Server.Reset(ctx.Context, server)
	})
}

// RescueMode sets the rescue mode for the Hetzner cloud microOS server
func RescueMode(ctx *cluster.Cluster, sshKey *hcloud.SSHKey, server *hcloud.Server) (string, error) {
	var rootPassword string
	err := performServerAction(ctx, server, ActionRescueMode, func() (*hcloud.Action, *hcloud.Response, error) {
		res, resp, err := ctx.CloudClient.Server.EnableRescue(ctx.Context, server, hcloud.ServerEnableRescueOpts{
			Type:    hcloud.ServerRescueTypeLinux64,
			SSHKeys: []*hcloud.SSHKey{sshKey},
		})
		if err == nil {
			rootPassword = res.RootPassword
		}
		return res.Action, resp, err
	})
	if err != nil {
		return "", err
	}

	if err := Reboot(ctx, server); err != nil {
		return "", err
	}

	return rootPassword, nil
}
