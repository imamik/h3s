package components

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
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
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
spec:
  ingressClassName: traefik
  rules:
    - host: {{ .K3sDomain }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: kubernetes
                port: 
                  number: 443
`,
		map[string]interface{}{
			"K3sDomain": fmt.Sprintf("k3s.%s", ctx.Config.Domain),
			"Host":      fmt.Sprintf("\\`k3s.%s\\`", ctx.Config.Domain),
		})
}

func K3sAPIServerConfig(ctx clustercontext.ClusterContext) string {
	return fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
}
