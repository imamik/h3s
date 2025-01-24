package install

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"

	"github.com/spf13/cobra"
)

func runInstallK3s(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "InstallK3s", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "k3s",
			"operation": "install",
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
		"k3sVersion":      ctx.Config.K3sVersion,
		"controlPlaneNum": ctx.Config.ControlPlane.Pool.Nodes,
		"workerNum":       len(ctx.Config.WorkerPools),
	})

	l.AddEvent(logger.Info, "Installing K3s")
	if err := deps.InstallK3s(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err).
			WithOperation("InstallK3s").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("k3sVersion", ctx.Config.K3sVersion).
			WithContext("controlPlaneNum", ctx.Config.ControlPlane.Pool.Nodes).
			WithContext("workerNum", len(ctx.Config.WorkerPools))
	}

	l.AddEvent(logger.Success, "K3s installed successfully")
	return nil
}

func runInstallComponents(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "InstallComponents", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "k8s-components",
			"operation": "install",
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
		"domain":     ctx.Config.Domain,
		"k3sVersion": ctx.Config.K3sVersion,
	})

	l.AddEvent(logger.Info, "Installing Kubernetes components")
	if err := deps.InstallK8sComponents(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err).
			WithOperation("InstallK8sComponents").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("k3sVersion", ctx.Config.K3sVersion)
	}

	l.AddEvent(logger.Success, "Kubernetes components installed successfully")
	return nil
}
