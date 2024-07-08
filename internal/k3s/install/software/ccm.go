package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ssh"
	"hcloud-k3s-cli/internal/utils/template"
)

func applySecretsCommand(ctx clustercontext.ClusterContext, network *hcloud.Network) string {
	return template.CompileTemplate(`kubectl apply -f - <<-EOF
apiVersion: "v1"
kind: "Secret"
metadata:
  namespace: 'kube-system'
  name: 'hcloud'
stringData:
  network: "{{ .NetworkName }}"
  token: "{{ .HCloudToken }}"
EOF`, map[string]interface{}{
		"NetworkName": network.Name,
		"HCloudToken": ctx.Credentials.HCloudToken,
	})
}

const applyInstallationsCommand = "kubectl apply -f https://github.com/hetznercloud/hcloud-cloud-controller-manager/releases/latest/download/ccm-networks.yaml"

func InstallHetznerCCM(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	ssh.ExecuteViaProxy(ctx, proxy, remote, applySecretsCommand(ctx, network))
	ssh.ExecuteViaProxy(ctx, proxy, remote, applyInstallationsCommand)
}
