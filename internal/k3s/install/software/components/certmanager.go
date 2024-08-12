package components

import (
	"h3s/internal/cluster"
)

const (
	CertManagerNamespace      = "cert-manager"
	CertManagerVersion        = "v1.15.1"
	CertManagerHetznerVersion = "1.3.1"
)

func CertManagerHelmChart() string {
	return kubectlApply(`
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager
  namespace: kube-system
spec:
  chart: cert-manager
  version: {{ .Version }}
  repo: https://charts.jetstack.io
  targetNamespace: {{ .Namespace }}
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
			"Version":   CertManagerVersion,
			"Namespace": CertManagerNamespace,
		})
}

func WaitForCertManagerCRDs() string {
	return WaitForCRDs("Cert-Manager", []string{
		"crd/certificaterequests.cert-manager.io",
		"crd/certificates.cert-manager.io",
		"crd/challenges.acme.cert-manager.io",
		"crd/clusterissuers.cert-manager.io",
		"crd/issuers.cert-manager.io",
		"crd/orders.acme.cert-manager.io",
	})
}

func CertManagerHetznerHelmChart(ctx *cluster.Cluster) string {
	return kubectlApply(`
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager-webhook-hetzner
  namespace: kube-system
spec:
  chart: cert-manager-webhook-hetzner
  version: {{ .Version }}
  repo: https://vadimkim.github.io/cert-manager-webhook-hetzner
  targetNamespace: {{ .Namespace }}
  set:
    groupName: {{ .Domain }}
`,
		map[string]interface{}{
			"Version":   CertManagerHetznerVersion,
			"Namespace": CertManagerNamespace,
			"Domain":    ctx.Config.Domain,
		})
}
