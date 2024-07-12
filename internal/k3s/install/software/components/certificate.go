package components

import (
	"encoding/base64"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
	"strings"
)

func WildcardCertificate(ctx clustercontext.ClusterContext) string {
	env := "staging"
	server := "https://acme-staging-v02.api.letsencrypt.org/directory"
	hetznerDNSTokenBase64 := base64.StdEncoding.EncodeToString([]byte(ctx.Credentials.HetznerDNSToken))

	if !ctx.Config.CertManager.Staging {
		env = "production"
		server = "https://acme-v02.api.letsencrypt.org/directory"
	}

	return template.CompileTemplate(`
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
  name: wildcard-issuer
  namespace: cert-manager
spec:
  acme:
    server: {{ .Server }}
    email: {{ .Email }}
    privateKeySecretRef:
      name: {{ .DomainKebap }}-{{ .Environment }}-issuer
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
  name: {{ .DomainKebap }}-wildcard
  namespace: cert-manager
spec:
  commonName: {{ .Domain }}
  secretName: {{ .DomainKebap }}-wildcard-tls
  issuerRef:
    name: {{ .DomainKebap }}-wildcard-http
    kind: ClusterIssuer
  dnsNames:
    - {{ .Domain }}
    - "*.{{ .Domain }}"
`,
		map[string]interface{}{
			"Environment":     env,
			"Server":          server,
			"Email":           ctx.Config.CertManager.Email,
			"HetznerDNSToken": hetznerDNSTokenBase64,
			"Domain":          ctx.Config.Domain,
			"DomainKebap":     strings.ReplaceAll(ctx.Config.Domain, ".", "-"),
		})
}