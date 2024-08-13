package components

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/k8s/components/templates"
	"h3s/internal/utils/template"
	"strings"
)

func CCM(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.CCM, getVars(ctx, network, lb))
}

func WildcardCertificate(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.Certificate, getVars(ctx, network, lb))
}

func CSIHelmChart(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.CSI, getVars(ctx, network, lb))
}

func CertManagerHelmChart(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.CertManager, getVars(ctx, network, lb))
}

func CertManagerHetznerHelmChart(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.CertManagerHetzner, getVars(ctx, network, lb))
}

func HCloudSecrets(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.HcloudSecrets, getVars(ctx, network, lb))
}

func K8sDashboardHelmChart(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.K8sDashboard, getVars(ctx, network, lb))
}

func K8sDashboardAccess(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.K8sDashboardConfig, getVars(ctx, network, lb))
}

func K8sIngress(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.K8sIngress, getVars(ctx, network, lb))
}

func TraefikHelmChart(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.Traefik, getVars(ctx, network, lb))
}

func TraefikDashboard(ctx *cluster.Cluster, network *hcloud.Network, lb *hcloud.LoadBalancer) string {
	return kubectlApply(templates.Yaml.TraefikDashboard, getVars(ctx, network, lb))
}

func K3sAPIServerConfig(ctx *cluster.Cluster) string {
	return fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
}

func kubectlApply(tpl string, data map[string]interface{}) string {
	yaml := template.CompileTemplate(tpl, data)
	cmd := strings.TrimSpace(yaml)
	fmt.Printf("\n%s\n", cmd)
	return "kubectl apply -f - <<EOF\n" + cmd + "\nEOF"
}
