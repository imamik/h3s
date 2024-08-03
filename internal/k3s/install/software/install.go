package software

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install/software/components"
	"h3s/internal/utils/ssh"
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

	ssh.ExecuteLocal(components.K3sAPIServerConfig(ctx))

}
