package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

const (
	microOsX86Link = "https://download.opensuse.org/tumbleweed/appliances/openSUSE-MicroOS.x86_64-ContainerHost-OpenStack-Cloud.qcow2"
	microOsArmLink = "https://download.opensuse.org/ports/aarch64/tumbleweed/appliances/openSUSE-MicroOS.aarch64-ContainerHost-OpenStack-Cloud.qcow2"
)

func Create(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
) *hcloud.Image {
	// log.Fatalf("createImage not implemented")
	return nil
}
