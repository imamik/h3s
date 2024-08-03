package components

import (
	"fmt"
	"h3s/internal/clustercontext"
)

func K3sAPI(ctx clustercontext.ClusterContext) string {
	return kubectlApply(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kubernetes-api-ingress
  namespace: default
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
    traefik.ingress.kubernetes.io/service.serverscheme: "https"
spec:
  ingressClassName: traefik
  rules:
    - host: {{ .Domain }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: kubernetes
                port: 
                  number: 443
  tls:
    - hosts:
        - {{ .Domain }}
      secretName: {{ .WildcardTLS }}
`,
		map[string]interface{}{
			"Domain":      fmt.Sprintf("k3s.%s", ctx.Config.Domain),
			"WildcardTLS": wildcardTlS(ctx),
		})
}

func K3sAPIServerConfig(ctx clustercontext.ClusterContext) string {
	return fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
}
