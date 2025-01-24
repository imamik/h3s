// Package components contains the functionality for managing the Kubernetes components & their configuration/variables
package components

import (
	"fmt"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/utils/encode"
	"regexp"
	"strings"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const (
	// CertManagerNamespace is the namespace where cert-manager is installed
	CertManagerNamespace = "cert-manager"
	// CertManagerVersion is the helm-version of cert-manager to install
	CertManagerVersion = "v1.15.1"
	// CertManagerHetznerVersion is the helm-version of the cert-manager-hetzner-dns01-solver to install
	CertManagerHetznerVersion = "1.3.1"

	// K8sDashboardNamespace is the namespace where the kubernetes-dashboard is installed
	K8sDashboardNamespace = "kubernetes-dashboard"
	// K8sDashboardVersion is the helm-version of the kubernetes-dashboard to install
	K8sDashboardVersion = "3.0.0"

	// TraefikNamespace is the namespace where traefik is installed
	TraefikNamespace = "traefik"
	// TraefikVersion is the helm-version of traefik to install
	TraefikVersion = "29.0.0"
	// TraefikImageTag is the docker-image-tag of traefik to install
	TraefikImageTag = "v3.1"
)

var (
	// TraefikCrds is a list of all traefik crds
	TraefikCrds = []string{
		"crd/accesscontrolpolicies.hub.traefik.io",
		"crd/apiaccesses.hub.traefik.io",
		"crd/apiportals.hub.traefik.io",
		"crd/apiratelimits.hub.traefik.io",
		"crd/apis.hub.traefik.io",
		"crd/apiversions.hub.traefik.io",
		"crd/ingressroutes.traefik.io",
		"crd/ingressroutetcps.traefik.io",
		"crd/ingressrouteudps.traefik.io",
		"crd/middlewares.traefik.io",
		"crd/middlewaretcps.traefik.io",
		"crd/serverstransports.traefik.io",
		"crd/serverstransporttcps.traefik.io",
		"crd/tlsoptions.traefik.io",
		"crd/tlsstores.traefik.io",
		"crd/traefikservices.traefik.io",
	}
	// CertManagerCrds is a list of all cert-manager crds
	CertManagerCrds = []string{
		"crd/certificaterequests.cert-manager.io",
		"crd/certificates.cert-manager.io",
		"crd/challenges.acme.cert-manager.io",
		"crd/clusterissuers.cert-manager.io",
		"crd/issuers.cert-manager.io",
		"crd/orders.acme.cert-manager.io",
	}
)

// kebapString converts a list of strings to a kebap-case string by joining them with a "-" (and removing all non-alphanumeric characters)
func kebapString(parts ...string) string {
	var res []string
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for _, s := range parts {
		parts := re.Split(s, -1)
		res = append(res, parts...)
	}
	return strings.Join(res, "-")
}

// GetVars returns a map of variables that are used in the yaml-templates
func GetVars(
	conf *config.Config,
	creds *credentials.ProjectCredentials,
	network *hcloud.Network,
	lb *hcloud.LoadBalancer,
) map[string]interface{} {
	domain := conf.Domain

	env := "staging"
	server := "https://acme-staging-v02.api.letsencrypt.org/directory"

	if conf.CertManager.Production {
		env = "production"
		server = "https://acme-v02.api.letsencrypt.org/directory"
	}

	return map[string]interface{}{
		"LoadbalancerName":     lb.Name,
		"LoadbalancerLocation": lb.Location.Name,

		"NetworkName":     network.Name,
		"ClusterCidrIpv4": network.IPRange.String(),
		"Server":          server,

		"Email":               conf.CertManager.Email,
		"PrivateKeySecretRef": kebapString(domain, env, "issuer"),

		"HetznerDNSToken": encode.ToBase64(creds.HetznerDNSToken),
		"HCloudToken":     creds.HCloudToken,

		"Domain": domain,

		"WildcardTLS":    kebapString(domain, "wildcard", "tls"),
		"WildcardIssuer": kebapString(domain, "wildcard", "issuer"),

		"CertManagerVersion":        CertManagerVersion,
		"CertManagerNamespace":      CertManagerNamespace,
		"CertManagerHetznerVersion": CertManagerHetznerVersion,

		"K8sDashboardVersion":   K8sDashboardVersion,
		"K8sDashboardNamespace": K8sDashboardNamespace,

		"TraefikNamespace":    TraefikNamespace,
		"TraefikVersion":      TraefikVersion,
		"TraefikImageTag":     TraefikImageTag,
		"TraefikHost":         fmt.Sprintf("\\`traefik.%s\\`", domain),
		"TraefikReplicaCount": 1,
	}
}
