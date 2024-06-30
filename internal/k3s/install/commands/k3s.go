package commands

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/template"
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

func K3sCommand(ctx clustercontext.ClusterContext) string {
	tplArr := []string{
		"curl -sfL https://get.k3s.io |",
		"INSTALL_K3S_SKIP_START=true",
		"INSTALL_K3S_SKIP_SELINUX_RPM=true",
		"INSTALL_K3S_CHANNEL={{ .InitialK3sChannel }}",
		"INSTALL_K3S_EXEC='server {{ .K3sExecServerArgs }}'",
		"sh -",
	}
	tpl := strings.Join(tplArr, " ")

	k3sExecServerArgs := ""
	k3sChannel := getMinorVersion(ctx.Config.K3sVersion)

	return template.CompileTemplate(tpl, map[string]interface{}{
		"InitialK3sChannel": k3sChannel,
		"K3sExecServerArgs": k3sExecServerArgs,
	})
}
