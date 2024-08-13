package k8s

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/server"
	"h3s/internal/k8s/components"
	"h3s/internal/utils/ssh"
	"h3s/internal/utils/template"
	"strings"
)

func Install(ctx *cluster.Cluster) error {
	net := network.Get(ctx)
	lb := loadbalancers.Get(ctx)
	gatewayNode, _ := gateway.GetIfNeeded(ctx)
	nodes, err := server.GetAll(ctx)
	if err != nil {
		return err
	}
	firstControlPlane := nodes.ControlPlane[0]
	_, err = installComponents(ctx, net, lb, gatewayNode, firstControlPlane)
	if err != nil {
		return err
	}
	return nil
}

func kubectlApply(tpl string, data map[string]interface{}) string {
	yaml := template.CompileTemplate(tpl, data)
	cmd := strings.TrimSpace(yaml)
	return "kubectl apply -f - <<EOF\n" + cmd + "\nEOF"
}

func waitForCRDs(component string, resources []string) string {
	waitCmd := "kubectl wait --for=condition=established --timeout=30s " + strings.Join(resources, " ") + " >/dev/null 2>&1"
	return fmt.Sprintf(`
echo "Waiting for CRDs of %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if %s; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, component, waitCmd)
}

func waitForNamespace(namespace string) string {
	return fmt.Sprintf(`
echo "Waiting for namespace %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if kubectl get namespace %s >/dev/null 2>&1; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, namespace, namespace)
}

func installComponents(
	ctx *cluster.Cluster,
	net *hcloud.Network,
	lb *hcloud.LoadBalancer,
	gateway *hcloud.Server,
	remote *hcloud.Server,
) (string, error) {
	vars := components.GetVars(ctx, net, lb)

	cmdArr := []string{
		// Setup Secrets
		kubectlApply(components.Yaml.HcloudSecrets, vars),

		// Install Hetzner CCM (Cloud Controller Manager)
		kubectlApply(components.Yaml.CCM, vars),

		// Install Hetzner CSI (Cloud Storage Interface)
		kubectlApply(components.Yaml.CSI, vars),

		// Install Cert-Manager
		kubectlApply(components.Yaml.CertManager, vars),
		waitForCRDs("Cert-Manager", components.CertManagerCrds),
		kubectlApply(components.Yaml.CertManagerHetzner, vars),

		// Install Traefik
		kubectlApply(components.Yaml.Traefik, vars),
		waitForCRDs("Traefik", components.TraefikCrds),

		// Install Wildcard Certificate
		kubectlApply(components.Yaml.Certificate, vars),

		// Setup K8s Dashboard
		kubectlApply(components.Yaml.K8sDashboard, vars),
		waitForNamespace(components.K8sDashboardNamespace),
		kubectlApply(components.Yaml.K8sDashboardConfig, vars),

		// Setup Traefik Dashboard
		kubectlApply(components.Yaml.TraefikDashboard, vars),

		// Configure K3s API Server Endpoint
		kubectlApply(components.Yaml.K8sIngress, vars),
	}

	for _, cmd := range cmdArr {
		_, err := ssh.ExecuteViaProxy(ctx, gateway, remote, cmd)
		if err != nil {
			fmt.Printf("Failed to apply: %s", err.Error())
		}
	}

	// Configure K3s API Server Endpoint
	clusterCmd := fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
	return ssh.ExecuteLocal(clusterCmd)

}
