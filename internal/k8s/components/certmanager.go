package components

import (
	_ "embed"
	"h3s/internal/cluster"
)

const (
	CertManagerNamespace      = "cert-manager"
	CertManagerVersion        = "v1.15.1"
	CertManagerHetznerVersion = "1.3.1"
)

//go:embed certmanager.yaml
var certManagerYAML string

//go:embed certmanager-hetzner.yaml
var certManagerHetznerYAML string

func CertManagerHelmChart() string {
	return kubectlApply(certManagerYAML, map[string]interface{}{
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
	return kubectlApply(certManagerHetznerYAML, map[string]interface{}{
		"Version":   CertManagerHetznerVersion,
		"Namespace": CertManagerNamespace,
		"Domain":    ctx.Config.Domain,
	})
}
