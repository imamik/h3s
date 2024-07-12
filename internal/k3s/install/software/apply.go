package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

func apply(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
	content string,
) {
	content = strings.TrimSpace(content)
	yaml := "- <<EOF\n" + content + "\nEOF"
	cmd := "kubectl apply -f " + yaml
	ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
}
