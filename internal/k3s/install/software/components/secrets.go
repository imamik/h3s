package components

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
)

func HCloudSecrets(ctx clustercontext.ClusterContext, network *hcloud.Network) string {
	return template.CompileTemplate(`
apiVersion: "v1"
kind: "Secret"
metadata:
  namespace: 'kube-system'
  name: 'hcloud'
stringData:
  network: "{{ .NetworkName }}"
  token: "{{ .HCloudToken }}"
`, map[string]interface{}{
		"NetworkName": network.Name,
		"HCloudToken": ctx.Credentials.HCloudToken,
	})
}
