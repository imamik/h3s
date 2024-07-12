package components

import (
	"hcloud-k3s-cli/internal/utils/template"
)

const (
	CertManagerNamespace = "cert-manager"
	CertManagerVersion   = "v1.15.1"
)

func CertManagerHelmChart() string {
	return template.CompileTemplate(`
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager
  namespace: kube-system
spec:
  chart: cert-manager
  version: {{ .Version }}
  repo: https://charts.jetstack.io
  targetNamespace: {{ .TargetNamespace }}
  createNamespace: true
  valuesContent: |-
    crds:
      enabled: true
    webhook:
      enabled: true
    cainjector:
      enabled: true
    startupapicheck:
      enabled: true
    ingressShim:
      enabled: true
`,
		map[string]interface{}{
			"Version":         CertManagerVersion,
			"TargetNamespace": CertManagerNamespace,
		})
}
