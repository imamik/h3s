package components

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/template"
	"strings"
)

const (
	TraefikNamespace = "traefik"
	TraefikVersion   = "29.0.0"
	TraefikImageTag  = "v3.1"
)

func valuesContent(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) string {
	values := template.CompileTemplate(`
image:
  tag: {{ .TraefikImageTag }}
deployment:
  replicas: {{ .ReplicaCount }}
globalArguments:
  - "--serversTransport.insecureSkipVerify=true"
service:
  enabled: true
  type: LoadBalancer
  annotations:
    "load-balancer.hetzner.cloud/name": "{{ .LoadbalancerName }}"
    "load-balancer.hetzner.cloud/use-private-ip": "true"
    "load-balancer.hetzner.cloud/disable-private-ingress": "false"
    "load-balancer.hetzner.cloud/disable-public-network": "false"
    "load-balancer.hetzner.cloud/ipv6-disabled": "false"
    "load-balancer.hetzner.cloud/location": "{{ .LoadbalancerLocation }}"
    "load-balancer.hetzner.cloud/type": "lb11"
    "load-balancer.hetzner.cloud/uses-proxyprotocol": "true"
    "load-balancer.hetzner.cloud/algorithm-type": "round_robin"
    "load-balancer.hetzner.cloud/health-check-interval": "5s"
    "load-balancer.hetzner.cloud/health-check-timeout": "3s"
    "load-balancer.hetzner.cloud/health-check-retries": "3"
ports:
  web:
    redirectTo:
      port: websecure
    proxyProtocol:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
    forwardedHeaders:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
  websecure:
    proxyProtocol:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
    forwardedHeaders:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
additionalArguments:
  - "--providers.kubernetesingress.ingressendpoint.publishedservice={{ .Namespace }}/traefik"
  - "--api.dashboard=true"
`,
		map[string]interface{}{
			"TraefikImageTag":      TraefikImageTag,
			"ReplicaCount":         1,
			"LoadbalancerName":     lb.Name,
			"LoadbalancerLocation": ctx.Config.ControlPlane.Pool.Location,
			"Namespace":            TraefikNamespace,
		})
	lines := strings.Split(values, "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	values = strings.Join(lines[1:], "\n")
	return values
}

func TraefikHelmChart(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) string {
	return kubectlApply(`
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Namespace }}
---
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: traefik
  namespace: kube-system
spec:
  chart: traefik
  version: {{ .TraefikVersion }}
  repo: https://traefik.github.io/charts
  targetNamespace: {{ .Namespace }}
  createNamespace: true
  bootstrap: true
  valuesContent: |-
{{ .ValuesContent }}
`,
		map[string]interface{}{
			"Namespace":      TraefikNamespace,
			"TraefikVersion": TraefikVersion,
			"ValuesContent":  valuesContent(ctx, lb),
		})
}

func WaitForTraefikCRDs() string {
	return WaitForCRDs("Traefik", []string{
		"crd/accesscontrolpolicies.hub.traefik.io",
		"crd/apiaccesses.hub.traefik.io",
		"crd/apiportals.hub.traefik.io",
		"crd/apiratelimits.hub.traefik.io",
		"crd/apis.hub.traefik.io",
		"crd/apiversions.hub.traefik.io",
		"crd/ingressroutes.traefik.io",
		"crd/ingressroutetcps.traefik.io",
		"crd/ingressrouteudps.traefik.io",
		"crd/middlewares.traefik.io",
		"crd/middlewaretcps.traefik.io",
		"crd/serverstransports.traefik.io",
		"crd/serverstransporttcps.traefik.io",
		"crd/tlsoptions.traefik.io",
		"crd/tlsstores.traefik.io",
		"crd/traefikservices.traefik.io",
	})
}

func TraefikDashboard(ctx clustercontext.ClusterContext) string {
	return kubectlApply(`
apiVersion: v1
kind: Secret
metadata:
  name: traefik-dashboard-auth-secret
  namespace: {{ .Namespace }}
type: kubernetes.io/basic-auth
stringData:
  username: admin
  password: pass
---
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: traefik-dashboard-auth
  namespace: {{ .Namespace }}
spec:
  basicAuth:
    secret: traefik-dashboard-auth-secret
---
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: traefik-dashboard
  namespace: {{ .Namespace }}
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  entryPoints:
    - websecure
  routes:
  - match: Host({{ .Host }})
    kind: Rule
    services:
    - name: api@internal
      kind: TraefikService
    middlewares:
    - name: traefik-dashboard-auth
`,
		map[string]interface{}{
			"Namespace": TraefikNamespace,
			"Host":      fmt.Sprintf("\\`traefik.%s\\`", ctx.Config.Domain),
		})
}
