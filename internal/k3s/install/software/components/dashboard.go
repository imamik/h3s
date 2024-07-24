package components

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"strings"
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
	return WaitForNamespaceToBeEstablished(K8sDashboardNamespace)
}

func K8sDashboardIngress(ctx clustercontext.ClusterContext) string {
	return kubectlApply(`
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
			"WildcardTLS": strings.ReplaceAll(ctx.Config.Domain, ".", "-") + "-wildcard-tls",
		})
}
