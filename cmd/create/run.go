package create

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"

	"github.com/spf13/cobra"
)

func runCreateConfig(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "CreateConfig", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "config",
			"operation": "create",
		})
	defer l.LogEvents()

	l.AddEvent(logger.Initialized)

	k3sReleases, err := deps.GetK3sReleases(false, true, 5)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to get k3s releases", err).
			WithOperation("GetK3sReleases").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("fetchStable", true).
			WithContext("fetchRC", false).
			WithContext("limit", 5)
	}

	if err := deps.BuildClusterConfig(k3sReleases); err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to build configuration", err).
			WithOperation("BuildClusterConfig").
			WithSeverity(errors.SeverityError).
			WithContext("releasesCount", len(k3sReleases))
	}

	l.AddEvent(logger.Success, "Configuration created successfully")
	return nil
}

func runCreateCredentials(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "CreateCredentials", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "credentials",
			"operation": "create",
		})
	defer l.LogEvents()

	l.AddEvent(logger.Initialized)

	if err := deps.ConfigureCredentials(); err != nil {
		return errors.Wrap(errors.ErrorTypeAuth, "failed to configure credentials", err).
			WithOperation("ConfigureCredentials").
			WithSeverity(errors.SeverityError).
			WithContext("component", "credentials")
	}

	l.AddEvent(logger.Success, "Credentials created successfully")
	return nil
}

func runCreateCluster(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "CreateCluster", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "cluster",
			"operation": "create",
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

	// Create Hetzner resources
	l.AddEvent(logger.Info, "Creating Hetzner resources")
	if err := deps.CreateHetznerResources(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to create hetzner resources", err).
			WithOperation("CreateHetznerResources").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("location", ctx.Config.ControlPlane.Pool.Location).
			WithContext("controlPlaneNum", ctx.Config.ControlPlane.Pool.Nodes).
			WithContext("workerNum", len(ctx.Config.WorkerPools))
	}
	l.AddEvent(logger.Success, "Hetzner resources created successfully")

	// Install K3s
	l.AddEvent(logger.Info, "Installing K3s")
	if err := deps.InstallK3s(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err).
			WithOperation("InstallK3s").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("k3sVersion", ctx.Config.K3sVersion).
			WithContext("domain", ctx.Config.Domain)
	}
	l.AddEvent(logger.Success, "K3s installed successfully")

	// Install Kubernetes components
	l.AddEvent(logger.Info, "Installing Kubernetes components")
	if err := deps.InstallK8sComponents(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err).
			WithOperation("InstallK8sComponents").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("k3sVersion", ctx.Config.K3sVersion).
			WithContext("domain", ctx.Config.Domain)
	}
	l.AddEvent(logger.Success, "Kubernetes components installed successfully")

	// Download kubeconfig
	l.AddEvent(logger.Info, "Downloading kubeconfig")
	if err := deps.DownloadKubeconfig(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig file", err).
			WithOperation("DownloadKubeconfig").
			WithSeverity(errors.SeverityWarning).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain)
	}
	l.AddEvent(logger.Success, "Kubeconfig downloaded successfully")

	l.AddEvent(logger.Success, "Cluster created successfully")
	return nil
}
