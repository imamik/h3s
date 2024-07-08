package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
	"strings"
)

func valuesContent(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	namespace string,
	imageTag string,
) string {
	values := template.CompileTemplate(`
image:
  tag: {{ .TraefikImageTag }}
deployment:
  replicas: {{ .IngressReplicaCount }}
globalArguments: []
service:
  enabled: true
  type: LoadBalancer
  annotations:
    "load-balancer.hetzner.cloud/name": "{{ .LoadbalancerName }}"
    "load-balancer.hetzner.cloud/use-private-ip": "true"
    "load-balancer.hetzner.cloud/disable-private-ingress": "true"
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
  kubectl:
	port: 6443
    expose:
      default: true
    exposedPort: 6443
    protocol: TCP
    proxyProtocol:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
    forwardedHeaders:
      trustedIPs:
        - 127.0.0.1/32
        - 10.0.0.0/8
additionalArguments:
  - "--providers.kubernetesingress.ingressendpoint.publishedservice={{ .IngressControllerNamespace }}/traefik"
`,
		map[string]interface{}{
			"TraefikImageTag":            imageTag,
			"IngressReplicaCount":        1,
			"LoadbalancerName":           lb.Name,
			"LoadbalancerLocation":       ctx.Config.ControlPlane.Pool.Location,
			"IngressControllerNamespace": namespace,
		})
	lines := strings.Split(values, "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	// start with the second line
	values = strings.Join(lines[1:], "\n")
	return values
}

func traefikYaml(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) string {
	namspace := "kube-system"
	version := "29.0.0"
	imageTag := "v3.1"
	return template.CompileTemplate(`
---
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .TargetNamespace }}
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
  targetNamespace: {{ .TargetNamespace }}
  bootstrap: true
  valuesContent: |-
{{ .ValuesContent }}
`,
		map[string]interface{}{
			"TargetNamespace": namspace,
			"TraefikVersion":  version,
			"ValuesContent":   valuesContent(ctx, lb, namspace, imageTag),
		})
}

func InstallTraefik(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	ApplyDynamicFile(ctx, proxy, remote, traefikYaml(ctx, lb))
}
