package components

import "hcloud-k3s-cli/internal/utils/template"

const (
	CertManagerNamespace = "cert-manager"
	CertManagerVersion   = "v1.15.1"
)

func CertManagerHelmChart() string {
	return template.CompileTemplate(`
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .TargetNamespace }}
---
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager
  namespace: kube-system
spec:
  chart: jetstack/cert-manager
  version: {{ .Version }}
  repo: https://charts.jetstack.io
  targetNamespace: {{ .TargetNamespace }}
`,
		map[string]interface{}{
			"Version":         CertManagerVersion,
			"TargetNamespace": CertManagerNamespace,
		})
}
