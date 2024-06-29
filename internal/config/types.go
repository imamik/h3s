package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

type Location string

const (
	Nuernberg   Location = "nbg1" // DE - EU Central
	Falkenstein Location = "fsn1" // DE - EU Central
	Helsinki    Location = "hel1" // FI - EU Central
	Ashburn     Location = "ash"  // US - US East
	Hillsboro   Location = "hil"  // US - US West
)

type CloudInstanceType string

const (
	SharedArm    CloudInstanceType = "shared-arm"
	SharedIntel  CloudInstanceType = "shared-intel"
	SharedAmd    CloudInstanceType = "shared-amd"
	DedicatedAmd CloudInstanceType = "dedicated-amd"
)

type CloudInstance string

const (
	CPX11 CloudInstance = "cpx11"
	CPX21 CloudInstance = "cpx21"
	CPX31 CloudInstance = "cpx31"
	CPX41 CloudInstance = "cpx41"
	CPX51 CloudInstance = "cpx51"
	CX22  CloudInstance = "cx22"
	CX32  CloudInstance = "cx32"
	CX42  CloudInstance = "cx42"
	CX52  CloudInstance = "cx52"
	CAX11 CloudInstance = "cax11"
	CAX21 CloudInstance = "cax21"
	CAX31 CloudInstance = "cax31"
	CAX41 CloudInstance = "cax41"
	CCX13 CloudInstance = "ccx13"
	CCX23 CloudInstance = "ccx23"
	CCX33 CloudInstance = "ccx33"
	CCX43 CloudInstance = "ccx43"
	CCX53 CloudInstance = "ccx53"
	CCX63 CloudInstance = "ccx63"
)

func GetArchitecture(instance CloudInstance) hcloud.Architecture {
	if instance[:3] == "cax" {
		return hcloud.ArchitectureARM
	} else {
		return hcloud.ArchitectureX86
	}
}
