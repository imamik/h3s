package kubectl

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"
	"strings"

	"github.com/spf13/cobra"
)

func runKubectl(cmd *cobra.Command, args []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "Kubectl", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "kubectl",
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
		"domain": ctx.Config.Domain,
	})

	l.AddEvent(logger.Info, "Executing kubectl command")
	res, err := deps.ExecuteKubectlCommand(ctx, args)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeKubectl, "failed to execute kubectl command", err).
			WithOperation("ExecuteKubectlCommand").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("args", strings.Join(args, " "))
	}

	l.AddEvent(logger.Success, "Kubectl command executed successfully")
	cmd.Println(res)
	return nil
}
