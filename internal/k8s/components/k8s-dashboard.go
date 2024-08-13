package components

import (
	_ "embed"
	"h3s/internal/cluster"
)

const (
	K8sDashboardVersion   = "3.0.0"
	K8sDashboardNamespace = "kubernetes-dashboard"
)

//go:embed k8s-dashboard.yaml
var k8sDashboardYAML string

//go:embed k8s-dashboard-config.yaml
var k8sDashboardConfigYAML string

func K8sDashboardHelmChart() string {
	return kubectlApply(k8sDashboardYAML, map[string]interface{}{
		"Version":   K8sDashboardVersion,
		"Namespace": K8sDashboardNamespace,
	})
}

func WaitForK8sDashboardNamespace() string {
	return WaitForNamespace(K8sDashboardNamespace)
}

func K8sDashboardAccess(ctx *cluster.Cluster) string {
	return kubectlApply(k8sDashboardConfigYAML, map[string]interface{}{
		"Namespace":   K8sDashboardNamespace,
		"Domain":      ctx.Config.Domain,
		"WildcardTLS": wildcardTlS(ctx),
	})
}
