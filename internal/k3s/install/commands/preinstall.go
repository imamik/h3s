package commands

import (
	"h3s/internal/clustercontext"
	"h3s/internal/utils/template"
)

func PreInstallCommand(ctx clustercontext.ClusterContext, configYaml string) string {
	return template.CompileTemplate(`
set -ex

# prepare the k3s config directory
mkdir -p /etc/rancher/k3s

# move the config file into place and adjust permissions
echo "{{ .ConfigYaml }}" > /etc/rancher/k3s/config.yaml
chmod 0600 /etc/rancher/k3s/config.yaml

# if the server has already been initialized just stop here
[ -e /etc/rancher/k3s/k3s.yaml ] && exit 0

# install kubectl
cat > /etc/profile.d/00-alias.sh <<EOF
alias k=kubectl
EOF

# install bash completion for kubectl
cat > /etc/bash_completion.d/kubectl <<EOF
if command -v kubectl >/dev/null; then
	source <(kubectl completion bash)
	complete -o default -F __start_kubectl k
fi
EOF

# wait for the internet connection to be available
timeout 180s /bin/sh -c 'while ! ping -c 1 {{ .Address }} >/dev/null 2>&1; do echo \"Ready for k3s installation, waiting for a successful connection to the internet...\"; sleep 5; done; echo Connected'
`, map[string]interface{}{
		"ConfigYaml":         configYaml,
		"Address":            "cloudflare.com",
		"OnlyPrivateNetwork": ctx.Config.PublicIps == false,
	})
}
