package components

import (
	_ "embed"
	"encoding/base64"
	"h3s/internal/cluster"
	"strings"
)

//go:embed certificate.yaml
var certificateYAML string

func domainKebap(ctx *cluster.Cluster) string {
	return strings.ReplaceAll(ctx.Config.Domain, ".", "-")
}

func wildcardTlS(ctx *cluster.Cluster) string {
	return domainKebap(ctx) + "-wildcard-tls"
}

func wildcardIssuer(ctx *cluster.Cluster) string {
	return domainKebap(ctx) + "-wildcard-issuer"
}

func WildcardCertificate(ctx *cluster.Cluster) string {
	env := "staging"
	server := "https://acme-staging-v02.api.letsencrypt.org/directory"
	hetznerDNSTokenBase64 := base64.StdEncoding.EncodeToString([]byte(ctx.Credentials.HetznerDNSToken))

	if ctx.Config.CertManager.Production {
		env = "production"
		server = "https://acme-v02.api.letsencrypt.org/directory"
	}

	return kubectlApply(certificateYAML, map[string]interface{}{
		"Server":              server,
		"Email":               ctx.Config.CertManager.Email,
		"HetznerDNSToken":     hetznerDNSTokenBase64,
		"Domain":              ctx.Config.Domain,
		"WildcardTLS":         wildcardTlS(ctx),
		"WildcardIssuer":      wildcardIssuer(ctx),
		"PrivateKeySecretRef": domainKebap(ctx) + "-" + env + "-issuer",
	})
}
