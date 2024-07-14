package components

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
)

func K3sAPI(ctx clustercontext.ClusterContext) string {
	return kubectlApply(`
apiVersion: v1
kind: Service
metadata:
  name: kubernetes-api-proxy
  namespace: kube-system
spec:
  ports:
    - port: 443
      targetPort: 6443
---
apiVersion: v1
kind: Endpoints
metadata:
  name: kubernetes-api-proxy
  namespace: kube-system
subsets:
  - addresses:
      - ip: 10.0.0.3
    ports:
      - port: 6443
---
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: kubernetes-api-tls
  namespace: kube-system
spec:
  entryPoints:
    - websecure
  routes:
    - match: Host({{ .Host }})
      kind: Rule
      services:
        - name: kubernetes-api-proxy
          port: 443
  tls: {}
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
    - host: {{ .K3sDomain }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: kubernetes-api-proxy
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
