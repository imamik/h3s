package microos

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// ImageInArchitecture is a struct representing the image in an architecture
type ImageInArchitecture struct {
	ARM *hcloud.Image // ARM is the image for the ARM architecture
	X86 *hcloud.Image // X86 is the image for the X86 architecture
}
