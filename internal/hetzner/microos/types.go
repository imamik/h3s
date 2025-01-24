package microos

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// ImageInArchitecture represents MicroOS images for different architectures
type ImageInArchitecture struct {
	ARM *hcloud.Image
	X86 *hcloud.Image
}
