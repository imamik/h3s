package components

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
)

func HCloudSecrets(ctx clustercontext.ClusterContext, network *hcloud.Network) string {
	return kubectlApply(`
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
