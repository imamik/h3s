package components

import (
	_ "embed"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
)

const (
	TraefikNamespace = "traefik"
	TraefikVersion   = "29.0.0"
	TraefikImageTag  = "v3.1"
)

//go:embed traefik.yaml
var traefikYAML string

//go:embed traefik-dashboard.yaml
var traefikDashboard string

func TraefikHelmChart(
	ctx *cluster.Cluster,
	lb *hcloud.LoadBalancer,
) string {
	return kubectlApply(traefikYAML, map[string]interface{}{
		"Namespace":            TraefikNamespace,
		"TraefikVersion":       TraefikVersion,
		"TraefikImageTag":      TraefikImageTag,
		"ReplicaCount":         1,
		"LoadbalancerName":     lb.Name,
		"LoadbalancerLocation": ctx.Config.ControlPlane.Pool.Location,
	})
}

func TraefikDashboard(ctx *cluster.Cluster) string {
	return kubectlApply(traefikDashboard, map[string]interface{}{
		"Namespace": TraefikNamespace,
		"Host":      fmt.Sprintf("\\`traefik.%s\\`", ctx.Config.Domain),
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
