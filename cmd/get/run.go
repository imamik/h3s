package get

import (
	"fmt"
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"h3s/internal/utils/logger"

	"github.com/spf13/cobra"
)

func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "GetKubeConfig", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "kubeconfig",
			"operation": "get",
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

	l.AddEvent(logger.Info, "Downloading kubeconfig")
	if err := deps.DownloadKubeconfig(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig", err).
			WithOperation("DownloadKubeconfig").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain)
	}

	l.AddEvent(logger.Success, "Kubeconfig downloaded successfully")
	return nil
}

func runGetToken(cmd *cobra.Command, _ []string) error {
	deps := dependencies.Get()
	l := logger.New(nil, logger.Cluster, "GetToken", "").
		WithLevel(logger.LevelInfo).
		WithFields(logger.LogFields{
			"component": "token",
			"operation": "get",
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

	l.AddEvent(logger.Info, "Generating bearer token")
	token, err := deps.GenerateK8sToken(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to generate bearer token", err).
			WithOperation("GenerateK8sToken").
			WithSeverity(errors.SeverityError).
			WithRetryable(true).
			WithContext("domain", ctx.Config.Domain).
			WithContext("namespace", "kubernetes-dashboard").
			WithContext("serviceAccount", "admin-user").
			WithContext("expirationHours", 24)
	}

	l.AddEvent(logger.Info, "Copying bearer token to clipboard")
	localCmd := fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", token)
	if _, err := deps.ExecuteLocalCommand(localCmd); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to copy bearer token to clipboard", err).
			WithOperation("ExecuteLocalCommand").
			WithSeverity(errors.SeverityWarning).
			WithRetryable(true)
	}

	l.AddEvent(logger.Success, "Bearer token copied to clipboard")
	cmd.Println("Bearer token copied to clipboard.")
	return nil
}
