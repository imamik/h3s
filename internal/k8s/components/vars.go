package components

import (
	"encoding/base64"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"regexp"
	"strings"
)

const (
	CertManagerNamespace      = "cert-manager"
	CertManagerVersion        = "v1.15.1"
	CertManagerHetznerVersion = "1.3.1"

	K8sDashboardNamespace = "kubernetes-dashboard"
	K8sDashboardVersion   = "3.0.0"

	TraefikNamespace = "traefik"
	TraefikVersion   = "29.0.0"
	TraefikImageTag  = "v3.1"
)

func kebapString(parts ...string) string {
	var res []string
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for _, s := range parts {
		parts := re.Split(s, -1)
		res = append(res, parts...)
	}
	return strings.Join(res, "-")
}

func getVars(
	ctx *cluster.Cluster,
	network *hcloud.Network,
	lb *hcloud.LoadBalancer,
) map[string]interface{} {
	conf := ctx.Config
	domain := conf.Domain

	env := "staging"
	server := "https://acme-staging-v02.api.letsencrypt.org/directory"
	HetznerDNSToken := base64.StdEncoding.EncodeToString([]byte(ctx.Credentials.HetznerDNSToken))

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

		"Email": conf.CertManager.Email,

		"HetznerDNSToken": HetznerDNSToken,
		"HCloudToken":     ctx.Credentials.HCloudToken,

		"Domain": domain,

		"WildcardTLS":         kebapString(domain, "wildcard", "tls"),
		"WildcardIssuer":      kebapString(domain, "wildcard", "issuer"),
		"PrivateKeySecretRef": kebapString(domain, env, "issuer"),

		"CertManagerVersion":   CertManagerVersion,
		"CertManagerNamespace": CertManagerNamespace,

		"CertManagerHetznerVersion": CertManagerHetznerVersion,

		"K8sDashboardVersion":   K8sDashboardVersion,
		"K8sDashboardNamespace": K8sDashboardNamespace,

		"TraefikNamespace":    TraefikNamespace,
		"TraefikVersion":      TraefikVersion,
		"TraefikImageTag":     TraefikImageTag,
		"TraefikHost":         fmt.Sprintf("\\`traefik.%s\\`", ctx.Config.Domain),
		"TraefikReplicaCount": 1,
	}
}
