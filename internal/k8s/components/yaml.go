package components

import (
	_ "embed"
)

// YamlTemplates contains all the YAML templates for the components, embedded into the binary
type YamlTemplates struct {
	CCM                string
	Certificate        string
	CertManager        string
	CertManagerHetzner string
	CSI                string
	HcloudSecrets      string
	K8sDashboard       string
	K8sDashboardConfig string
	K8sIngress         string
	Traefik            string
	TraefikDashboard   string
}

//go:embed ccm.yaml
var ccmYAML string

//go:embed certificate.yaml
var certificateYAML string

//go:embed certmanager.yaml
var certManagerYAML string

//go:embed certmanager-hetzner.yaml
var certManagerHetznerYAML string

//go:embed csi.yaml
var csiYAML string

//go:embed hcloud-secrets.yaml
var hcloudSecretsYAML string

//go:embed k8s-dashboard.yaml
var k8sDashboardYAML string

//go:embed k8s-dashboard-config.yaml
var k8sDashboardConfigYAML string

//go:embed k8s-ingress.yaml
var k8sIngressYAML string

//go:embed traefik.yaml
var traefikYAML string

//go:embed traefik-dashboard.yaml
var traefikDashboard string

// Yaml contains all the YAML templates for the components, embedded into the binary
var Yaml = YamlTemplates{
	CCM:                ccmYAML,
	Certificate:        certificateYAML,
	CertManager:        certManagerYAML,
	CertManagerHetzner: certManagerHetznerYAML,
	CSI:                csiYAML,
	HcloudSecrets:      hcloudSecretsYAML,
	K8sDashboard:       k8sDashboardYAML,
	K8sDashboardConfig: k8sDashboardConfigYAML,
	K8sIngress:         k8sIngressYAML,
	Traefik:            traefikYAML,
	TraefikDashboard:   traefikDashboard,
}
