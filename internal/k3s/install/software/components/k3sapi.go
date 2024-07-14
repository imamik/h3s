package components

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
)

func K3sAPI(ctx clustercontext.ClusterContext) string {
	return kubectlApply(`
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: kubernetes-api-tls
  namespace: kube-system
spec:
  entryPoints:
    - websecure
  routes:
    - match: Host(`+"k3s.{{ .Domain }}"+`)
      kind: Rule
      services:
        - name: kubernetes-api
          port: 6443
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kubernetes-api
  namespace: kube-system
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
spec:
  ingressClassName: traefik
  rules:
    - host: k3s.{{ .Domain }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: kubernetes-api
                port: 
                  number: 6443
`,
		map[string]interface{}{
			"Domain": ctx.Config.Domain,
		})
}

func K3sAPIServerConfig(ctx clustercontext.ClusterContext) string {
	return fmt.Sprintf("kubectl config set-cluster default --server=https://k3s.%s", ctx.Config.Domain)
}
