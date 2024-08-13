package components

import (
	_ "embed"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
)

//go:embed hcloud-secrets.yaml
var hcloudSecretsYAML string

func HCloudSecrets(ctx *cluster.Cluster, network *hcloud.Network) string {
	return kubectlApply(hcloudSecretsYAML, map[string]interface{}{
		"NetworkName": network.Name,
		"HCloudToken": ctx.Credentials.HCloudToken,
	})
}
