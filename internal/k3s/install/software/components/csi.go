package components

func CSIHelmChart() string {
	return `
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: hcloud-csi
  namespace: kube-system
spec:
  chart: hcloud/hcloud-csi
  version: 2.8.0
  repo: https://charts.hetzner.cloud
  targetNamespace: kube-system
`
}
