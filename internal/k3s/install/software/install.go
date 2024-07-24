package software

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/software/components"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func Install(
	ctx clustercontext.ClusterContext,
	net *hcloud.Network,
	lb *hcloud.LoadBalancer,
	gateway *hcloud.Server,
	remote *hcloud.Server,
) {
	cmdArr := []string{
		// Setup Secrets
		components.HCloudSecrets(ctx, net),

		// Install Hetzner CCM (Cloud Controller Manager)
		components.CCMServiceAccount(),
		components.CCMRoleBinding(),
		components.CCMSettings(ctx, net),

		// Install Hetzner CSI (Cloud Storage Interface)
		// components.CSIHelmChart(),

		// Install Cert-Manager
		components.CertManagerHelmChart(),
		components.WaitForCertManagerCRDs(),
		components.CertManagerHetznerHelmChart(ctx),

		// Install Traefik
		components.TraefikHelmChart(ctx, lb),
		components.WaitForTraefikCRDs(),

		// Install Wildcard Certificate
		components.WildcardCertificate(ctx),

		// Setup Dashboards
		components.K8sDashboardHelmChart(),
		components.K8sDashboardIngress(ctx),
		components.TraefikDashboard(ctx),

		// Configure K3s API Server Endpoint
		components.K3sAPI(ctx),
		// components.K3sAPIServerConfig(ctx),
	}

	for _, cmd := range cmdArr {
		_, err := ssh.ExecuteViaProxy(ctx, gateway, remote, cmd)
		if err != nil {
			fmt.Printf("Failed to apply: %s", err.Error())
		}
	}

}
