package components

import (
	_ "embed"
	"fmt"
	"h3s/internal/cluster"
)

//go:embed k8s-ingress.yaml
var k8sIngressYAML string

func K8sIngress(ctx *cluster.Cluster) string {
	return kubectlApply(k8sIngressYAML, map[string]interface{}{
		"Domain":      fmt.Sprintf("k3s.%s", ctx.Config.Domain),
		"WildcardTLS": wildcardTlS(ctx),
	})
}

func K3sAPIServerConfig(ctx *cluster.Cluster) string {
	return fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
}
