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

func installComponents(
	ctx *cluster.Cluster,
	net *hcloud.Network,
	lb *hcloud.LoadBalancer,
	gateway *hcloud.Server,
	remote *hcloud.Server,
) (string, error) {
	cmdArr := []string{
		// Setup Secrets
		components.HCloudSecrets(ctx, net),

		// Install Hetzner CCM (Cloud Controller Manager)
		components.CCMServiceAccount(),
		components.CCMRoleBinding(),
		components.CCMSettings(ctx, net),

		// Install Hetzner CSI (Cloud Storage Interface)
		components.CSIHelmChart(),

		// Install Cert-Manager
		components.CertManagerHelmChart(),
		components.WaitForCertManagerCRDs(),
		components.CertManagerHetznerHelmChart(ctx),

		// Install Traefik
		components.TraefikHelmChart(ctx, lb),
		components.WaitForTraefikCRDs(),

		// Install Wildcard Certificate
		components.WildcardCertificate(ctx),

		// Setup K8s Dashboard
		components.K8sDashboardHelmChart(),
		components.WaitForK8sDashboardNamespace(),
		components.K8sDashboardAccess(ctx),

		// Setup Traefik Dashboard
		components.TraefikDashboard(ctx),

		// Configure K3s API Server Endpoint
		components.K3sAPI(ctx),
	}

	for _, cmd := range cmdArr {
		_, err := ssh.ExecuteViaProxy(ctx, gateway, remote, cmd)
		if err != nil {
			fmt.Printf("Failed to apply: %s", err.Error())
		}
	}

	return ssh.ExecuteLocal(components.K3sAPIServerConfig(ctx))

}
