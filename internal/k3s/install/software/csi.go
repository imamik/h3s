package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

func hcloudCSIHelmChartYaml() string {
	return `
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: hcloud-csi-driver
  namespace: kube-system
spec:
  chart: hcloud/hcloud-csi-driver
  version: 2.8.0
  repo: https://charts.hetzner.cloud
  targetNamespace: kube-system
`
}

func InstallHetznerCSI(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	apply(ctx, proxy, remote, hcloudCSIHelmChartYaml())
}
