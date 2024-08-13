package components

import (
	_ "embed"
	"fmt"
	"strings"
)

func WaitForCertManagerCRDs() string {
	return waitForCRDs("Cert-Manager", []string{
		"crd/certificaterequests.cert-manager.io",
		"crd/certificates.cert-manager.io",
		"crd/challenges.acme.cert-manager.io",
		"crd/clusterissuers.cert-manager.io",
		"crd/issuers.cert-manager.io",
		"crd/orders.acme.cert-manager.io",
	})
}

func WaitForTraefikCRDs() string {
	return waitForCRDs("Traefik", []string{
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

func WaitForK8sDashboardNamespace() string {
	return waitForNamespace(K8sDashboardNamespace)
}

func waitForCRDs(component string, resources []string) string {
	waitCmd := "kubectl wait --for=condition=established --timeout=30s " + strings.Join(resources, " ") + " >/dev/null 2>&1"
	return fmt.Sprintf(`
echo "Waiting for CRDs of %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if %s; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, component, waitCmd)
}

func waitForNamespace(namespace string) string {
	return fmt.Sprintf(`
echo "Waiting for namespace %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if kubectl get namespace %s >/dev/null 2>&1; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, namespace, namespace)
}
