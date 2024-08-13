package commands

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

const (
	microOsX86Link = "https://download.opensuse.org/tumbleweed/appliances/openSUSE-MicroOS.x86_64-ContainerHost-OpenStack-Cloud.qcow2"
	microOsArmLink = "https://download.opensuse.org/ports/aarch64/tumbleweed/appliances/openSUSE-MicroOS.aarch64-ContainerHost-OpenStack-Cloud.qcow2"
	wget           = "wget --timeout=5 --waitretry=5 --tries=5 --retry-connrefused --inet4-only "
)

func DownloadImage(architecture hcloud.Architecture) string {
	if architecture == hcloud.ArchitectureARM {
		return wget + microOsArmLink
	}
	return wget + microOsX86Link
}
