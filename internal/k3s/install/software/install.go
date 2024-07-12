package software

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/software/components"
	"hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

func Install(
	ctx clustercontext.ClusterContext,
	net *hcloud.Network,
	lb *hcloud.LoadBalancer,
	gateway *hcloud.Server,
	remote *hcloud.Server,
) {

	apply := func(content string) {
		content = strings.TrimSpace(content)
		yaml := "- <<EOF\n" + content + "\nEOF"
		cmd := "kubectl apply -f " + yaml
		_, err := ssh.ExecuteViaProxy(ctx, gateway, remote, cmd)
		if err != nil {
			fmt.Printf("Failed to apply: %s", err.Error())
		}
	}

	// Setup Secrets
	apply(components.HCloudSecrets(ctx, net))

	// Install Hetzner CCM (Cloud Controller Manager)
	apply(components.CCMServiceAccount())
	apply(components.CCMRoleBinding())
	apply(components.CCMSettings(ctx, net))

	// Install Hetzner CSI (Cloud Storage Interface)
	apply(components.CSIHelmChart())

	// Install Traefik
	apply(components.TraefikNamespace())
	apply(components.TraefikHelmChartWithValues(ctx, lb))

	// Install Cert-Manager
	apply(components.CertManagerNamespace())
	apply(components.CertManagerHelmChart())

}
