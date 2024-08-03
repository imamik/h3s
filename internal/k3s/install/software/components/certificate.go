package components

import (
	"encoding/base64"
	"h3s/internal/clustercontext"
	"strings"
)

func domainKebap(ctx clustercontext.ClusterContext) string {
	return strings.ReplaceAll(ctx.Config.Domain, ".", "-")
}

func wildcardTlS(ctx clustercontext.ClusterContext) string {
	return domainKebap(ctx) + "-wildcard-tls"
}

func wildcardIssuer(ctx clustercontext.ClusterContext) string {
	return domainKebap(ctx) + "-wildcard-issuer"
}

func WildcardCertificate(ctx clustercontext.ClusterContext) string {
	env := "staging"
	server := "https://acme-staging-v02.api.letsencrypt.org/directory"
	hetznerDNSTokenBase64 := base64.StdEncoding.EncodeToString([]byte(ctx.Credentials.HetznerDNSToken))

	if ctx.Config.CertManager.Production {
		env = "production"
		server = "https://acme-v02.api.letsencrypt.org/directory"
	}

	return kubectlApply(`
apiVersion: v1
kind: Secret
metadata:
  name: hetzner-secret
  namespace: cert-manager
type: Opaque
data:
  api-key: {{ .HetznerDNSToken }}
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: {{ .WildcardIssuer }}
spec:
  acme:
    server: {{ .Server }}
    email: {{ .Email }}
    privateKeySecretRef:
      name: {{ .PrivateKeySecretRef }}
    solvers:
      - dns01:
          webhook:
            groupName: {{ .Domain }}
            solverName: hetzner
            config:
              secretName: hetzner-secret
              zoneName: {{ .Domain }}
              apiUrl: https://dns.hetzner.com/api/v1
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: {{ .WildcardTLS }}
spec:
  commonName: {{ .Domain }}
  secretName: {{ .WildcardTLS }}
  issuerRef:
    name: {{ .WildcardIssuer }}
    kind: ClusterIssuer
  dnsNames:
    - {{ .Domain }}
    - "*.{{ .Domain }}"
---
apiVersion: traefik.io/v1alpha1
kind: TLSStore
metadata:
  name: default
spec:
  defaultCertificate:
    secretName: {{ .WildcardTLS }}
`,
		map[string]interface{}{
			"Server":              server,
			"Email":               ctx.Config.CertManager.Email,
			"HetznerDNSToken":     hetznerDNSTokenBase64,
			"Domain":              ctx.Config.Domain,
			"WildcardTLS":         wildcardTlS(ctx),
			"WildcardIssuer":      wildcardIssuer(ctx),
			"PrivateKeySecretRef": domainKebap(ctx) + "-" + env + "-issuer",
		})
}
