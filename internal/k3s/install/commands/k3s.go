package commands

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
	"strings"
)

func getMinorVersion(version string) string {
	// Split the version string by "."
	versionParts := strings.Split(version, ".")

	// Check if there are at least two parts (major and minor version)
	if len(versionParts) >= 2 {
		// Concatenate the major and minor version
		return versionParts[0] + "." + versionParts[1]
	} else {
		// Handle the case where the version string is not in the expected format
		// This could be an error or a default value
		return "default"
	}
}

func K3sInstall(ctx clustercontext.ClusterContext, isControlPlane bool) string {
	tplArr := []string{
		"curl -sfL https://get.k3s.io | INSTALL_K3S_SKIP_START=true INSTALL_K3S_SKIP_SELINUX_RPM=true INSTALL_K3S_CHANNEL={{ .InitialK3sChannel }} INSTALL_K3S_EXEC='{{ .ServerOrAgent }} {{ .K3sExecServerArgs }}' sh -",
	}
	tpl := strings.Join(tplArr, " ")

	k3sExecArgs := ""
	k3sChannel := getMinorVersion(ctx.Config.K3sVersion)
	serverOrAgent := "agent"
	if isControlPlane {
		serverOrAgent = "server"
	}

	return template.CompileTemplate(tpl, map[string]interface{}{
		"InitialK3sChannel": k3sChannel,
		"K3sExecServerArgs": k3sExecArgs,
		"ServerOrAgent":     serverOrAgent,
	})
}

func K3sStartServer() string {
	return `
systemctl start k3s 2> /dev/null

# prepare the needed directories
mkdir -p /var/post_install /var/user_kustomize

# wait for the server to be ready
timeout 360 bash <<EOF
	until systemctl status k3s > /dev/null; do
		systemctl start k3s 2> /dev/null
		echo "Waiting for the k3s server to start..."
		sleep 3
	done
EOF
`
}

func K3sStartAgent() string {
	return `
systemctl start k3s-agent 2> /dev/null

# wait for the agent to be ready
timeout 120 bash <<EOF
	until systemctl status k3s-agent > /dev/null; do
		systemctl start k3s-agent 2> /dev/null
		echo "Waiting for the k3s agent to start..."
		sleep 2
	done
EOF
`
}
