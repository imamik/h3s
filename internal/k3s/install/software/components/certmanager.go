package components

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
	"strings"
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

func WaitForCRDsToBeEstablished() string {
	resources := []string{
		"crd/certificaterequests.cert-manager.io",
		"crd/certificates.cert-manager.io",
		"crd/challenges.acme.cert-manager.io",
		"crd/clusterissuers.cert-manager.io",
		"crd/issuers.cert-manager.io",
		"crd/orders.acme.cert-manager.io",
	}
	waitCmd := "kubectl wait --for=condition=established --timeout=30s " + strings.Join(resources, " ")
	return fmt.Sprintf(`
		for i in {1..5}; do
			if %s; then
				echo "CRDs established successfully"
				exit 0
			fi
			echo "Waiting for CRDs to be established (attempt $i)"
			sleep 10
		done
		echo "Timed out waiting for CRDs to be established"
		exit 1
	`, waitCmd)
}

func CertManagerHetznerHelmChart(ctx clustercontext.ClusterContext) string {
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
