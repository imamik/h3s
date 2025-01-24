package ssh

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"
	"strings"

	"github.com/spf13/cobra"
)

// runSSH proxies ssh commands to the first control plane server in the h3s cluster
func runSSH(cmd *cobra.Command, args []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Server, "SSH", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "ssh",
			"operation": strings.Join(args, " "),
		})
	defer l.LogEvents()

	l.AddEvent(logger.Initialized)

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err).
			WithOperation("GetClusterContext").
			WithSeverity(errors.SeverityFatal).
			WithContext("component", "cluster")
	}

	l.WithFields(logger.LogFields{
		"domain":          ctx.Config.Domain,
		"controlPlaneNum": ctx.Config.ControlPlane.Pool.Nodes,
	})

	command := strings.Join(args, " ")
	l.AddEvent(logger.Info, "Executing SSH command")
	res, err := deps.ExecuteSSHCommand(ctx, command)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeSSH, "failed to execute ssh command", err).
			WithOperation("ExecuteSSHCommand").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("command", command)
	}

	l.AddEvent(logger.Success, "SSH command executed successfully")
	cmd.Println(res)
	return nil
}
