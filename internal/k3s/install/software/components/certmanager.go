package components

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
)

const (
	CertManagerNamespace      = "cert-manager"
	CertManagerVersion        = "v1.15.1"
	CertManagerHetznerVersion = "1.3.1"
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

func CertManagerHetznerHelmChart(ctx clustercontext.ClusterContext) string {
	return template.CompileTemplate(`
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager-webhook-hetzner
  namespace: kube-system
spec:
  chart: cert-manager-webhook-hetzner
  version: {{ .Version }}
  repo: https://vadimkim.github.io/cert-manager-webhook-hetzner
  targetNamespace: {{ .TargetNamespace }}
  set:
    groupName: {{ .Domain }}
`,
		map[string]interface{}{
			"Version":         CertManagerHetznerVersion,
			"TargetNamespace": CertManagerNamespace,
			"Domain":          ctx.Config.Domain,
		})
}
