package destroy

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"

	"github.com/spf13/cobra"
)

func runDestroyCluster(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "DestroyCluster", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "cluster",
			"operation": "destroy",
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
		"location":        ctx.Config.ControlPlane.Pool.Location,
		"controlPlaneNum": ctx.Config.ControlPlane.Pool.Nodes,
		"workerNum":       len(ctx.Config.WorkerPools),
	})

	// Destroy Hetzner resources
	l.AddEvent(logger.Info, "Destroying Hetzner resources")
	if err := deps.DestroyHetznerResources(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to destroy cluster resources", err).
			WithOperation("DestroyHetznerResources").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("location", ctx.Config.ControlPlane.Pool.Location).
			WithContext("controlPlaneNum", ctx.Config.ControlPlane.Pool.Nodes).
			WithContext("workerNum", len(ctx.Config.WorkerPools))
	}
	l.AddEvent(logger.Success, "Hetzner resources destroyed successfully")

	l.AddEvent(logger.Success, "Cluster destroyed successfully")
	return nil
}
