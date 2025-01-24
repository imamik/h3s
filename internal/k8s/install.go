// Package k8s contains the functionality for installing the Kubernetes components
package k8s

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/server"
	"h3s/internal/k8s/components"
	"h3s/internal/utils/kubectl"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ssh"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Install installs the Kubernetes components to the first control plane node of the cluster
func Install(clr *cluster.Cluster) error {
	l := logger.New(nil, logger.Server, logger.Create, "gateway")
	defer l.LogEvents()

	// Get network
	net, err := network.Get(clr)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	// Get load balancer
	lb, err := loadbalancers.Get(clr)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	// Get gateway node
	gatewayNode, err := gateway.Get(clr)
	if err != nil {
		return fmt.Errorf("failed to get gateway node: %w", err)
	}

	// Get all nodes
	nodes, err := server.GetAll(clr)
	if err != nil {
		return fmt.Errorf("failed to get all nodes: %w", err)
	}

	firstControlPlane := nodes.ControlPlane[0]
	vars := components.GetVars(clr.Config, clr.Credentials, net, lb)

	return installComponents(clr, gatewayNode, firstControlPlane, vars)
}

func retryCommand(command, description string) string {
	return fmt.Sprintf(`
echo "Waiting for %s"
for i in {1..5}; do
	if %s; then
		[ "$i" -gt 1 ] && sleep 10
		echo "Successfully established"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1`, description, command)
}

func installComponents(clr *cluster.Cluster, gateway, remote *hcloud.Server, vars map[string]interface{}) error {
	steps := []struct {
		applyYaml        string
		waitForNamespace string
		description      string
		waitForCrds      []string
	}{
		{applyYaml: components.Yaml.HcloudSecrets, description: "Hetzner Cloud Secrets"},
		{applyYaml: components.Yaml.CCM, description: "Hetzner Cloud Controller Manager"},
		{applyYaml: components.Yaml.CSI, description: "Hetzner Cloud Storage Class"},
		{applyYaml: components.Yaml.CertManager, description: "Cert-Manager"},
		{waitForCrds: components.CertManagerCrds, description: "Cert-Manager CRDs"},
		{applyYaml: components.Yaml.CertManagerHetzner, description: "Cert-Manager Hetzner Issuer"},
		{applyYaml: components.Yaml.Traefik, description: "Traefik"},
		{waitForCrds: components.TraefikCrds, description: "Traefik CRDs"},
		{applyYaml: components.Yaml.Certificate, description: "Certificate"},
		{applyYaml: components.Yaml.K8sDashboard, description: "Kubernetes Dashboard"},
		{waitForNamespace: components.K8sDashboardNamespace, description: components.K8sDashboardNamespace},
		{applyYaml: components.Yaml.K8sDashboardConfig, description: "Kubernetes Dashboard Config"},
		{applyYaml: components.Yaml.TraefikDashboard, description: "Traefik Dashboard"},
		{applyYaml: components.Yaml.K8sIngress, description: "Kubernetes Ingress"},
	}

	for _, step := range steps {
		var kubectlCmd *kubectl.Command
		if step.applyYaml != "" {
			kubectlCmd = kubectl.New().ApplyTemplate(step.applyYaml, vars)
			step.description = fmt.Sprintf("Apply Component %s", step.description)
		}
		if step.waitForNamespace != "" {
			kubectlCmd = kubectl.New().GetResource("namespace " + step.waitForNamespace).DevNull()
			step.description = fmt.Sprintf("Wait for Namespace %s", step.description)
		}
		if step.waitForCrds != nil {
			kubectlCmd = kubectl.New().WaitForEstablished(step.waitForCrds...).DevNull()
			step.description = fmt.Sprintf("Wait for %v", step.description)
		}
		if kubectlCmd == nil {
			return fmt.Errorf("no command to execute")
		}
		cmd, err := kubectlCmd.String()
		if err != nil {
			return err
		}
		if step.waitForCrds != nil || step.waitForNamespace != "" {
			cmd = retryCommand(cmd, step.description)
		}
		if _, err := ssh.ExecuteViaProxy(clr.Config.SSHKeyPaths.PrivateKeyPath, gateway, remote, cmd); err != nil {
			return fmt.Errorf("command execution failed: %w", err)
		}
	}

	return nil
}
