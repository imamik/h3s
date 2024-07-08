package software

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func ApplyYaml(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
	yaml string,
) (string, error) {
	cmd := "kubectl apply -f " + yaml
	fmt.Println("Applying YAML: \n" + cmd)
	return ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
}

func ApplyDynamicFile(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
	content string,
) (string, error) {
	return ApplyYaml(ctx, proxy, remote, "- <<EOF\n"+content+"\nEOF")
}
