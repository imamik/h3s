package components

import (
	"h3s/internal/cluster"
)

const (
	K8sDashboardVersion   = "3.0.0"
	K8sDashboardNamespace = "kubernetes-dashboard"
)

func K8sDashboardHelmChart() string {
	return kubectlApply(`
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  chart: kubernetes-dashboard
  version: {{ .Version }}
  repo: https://kubernetes.github.io/dashboard/
  targetNamespace: {{ .Namespace }}
  createNamespace: true
`,
		map[string]interface{}{
			"Version":   K8sDashboardVersion,
			"Namespace": K8sDashboardNamespace,
		})
}

func WaitForK8sDashboardNamespace() string {
	return WaitForNamespace(K8sDashboardNamespace)
}

func K8sDashboardAccess(ctx *cluster.Cluster) string {
	return kubectlApply(`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kubernetes-dashboard-ingress
  namespace: {{ .Namespace }}
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  ingressClassName: traefik
  rules:
  - host: dashboard.{{ .Domain }}
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: kubernetes-dashboard
            port:
              number: 443
`,
		map[string]interface{}{
			"Namespace":   K8sDashboardNamespace,
			"Domain":      ctx.Config.Domain,
			"WildcardTLS": wildcardTlS(ctx),
		})
}
