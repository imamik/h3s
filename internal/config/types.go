package config

// Location is a type alias for a string representing a Hetzner location.
type Location string

const (
	Nuernberg   Location = "nbg1" // Nuernberg is a location in DE - EU Central
	Falkenstein Location = "fsn1" // Falkenstein is a location in DE - EU Central
	Helsinki    Location = "hel1" // Helsinki is a location in FI - EU Central
	Ashburn     Location = "ash"  // Ashburn is a location in US - US East
	Hillsboro   Location = "hil"  // Hillsboro is a location in US - US West
)

// CloudInstanceType is a type alias for a string representing a cloud instance type e.g. shared-intel, dedicated-amd, etc.
type CloudInstanceType string

const (
	SharedArm    CloudInstanceType = "shared-arm"    // SharedArm is an instance type, best value for general-purpose workloads
	SharedIntel  CloudInstanceType = "shared-intel"  // SharedIntel is an instance type, good for general-purpose workloads
	SharedAmd    CloudInstanceType = "shared-amd"    // SharedAmd is an instance type, good for general-purpose workloads
	DedicatedAmd CloudInstanceType = "dedicated-amd" // DedicatedAmd is an instance type, good for high-performance workloads
)

// CloudInstance is a type alias for a string representing a cloud instance e.g. cpx11, cx22, etc.
type CloudInstance string

const (
	CPX11 CloudInstance = "cpx11" // CPX11 is an instance with 2 vCPU, 2 GB RAM, Shared CPU (Intel)
	CPX21 CloudInstance = "cpx21" // CPX21 is an instance with 3 vCPU, 4 GB RAM, Shared CPU (Intel)
	CPX31 CloudInstance = "cpx31" // CPX31 is an instance with 4 vCPU, 8 GB RAM, Shared CPU (Intel)
	CPX41 CloudInstance = "cpx41" // CPX41 is an instance with 8 vCPU, 16 GB RAM, Shared CPU (Intel)
	CPX51 CloudInstance = "cpx51" // CPX51 is an instance with 16 vCPU, 32 GB RAM, Shared CPU (Intel)
	CX22  CloudInstance = "cx22"  // CX22 is an instance with 2 vCPU, 8 GB RAM, Shared CPU (AMD)
	CX32  CloudInstance = "cx32"  // CX32 is an instance with 4 vCPU, 16 GB RAM, Shared CPU (AMD)
	CX42  CloudInstance = "cx42"  // CX42 is an instance with 8 vCPU, 32 GB RAM, Shared CPU (AMD)
	CX52  CloudInstance = "cx52"  // CX52 is an instance with 16 vCPU, 64 GB RAM, Shared CPU (AMD)
	CAX11 CloudInstance = "cax11" // CAX11 is an instance with 2 vCPU, 4 GB RAM, Shared CPU (ARM)
	CAX21 CloudInstance = "cax21" // CAX21 is an instance with 4 vCPU, 8 GB RAM, Shared CPU (ARM)
	CAX31 CloudInstance = "cax31" // CAX31 is an instance with 8 vCPU, 16 GB RAM, Shared CPU (ARM)
	CAX41 CloudInstance = "cax41" // CAX41 is an instance with 16 vCPU, 32 GB RAM, Shared CPU (ARM)
	CCX13 CloudInstance = "ccx13" // CCX13 is an instance with 2 vCPU, 8 GB RAM, Dedicated CPU (AMD)
	CCX23 CloudInstance = "ccx23" // CCX23 is an instance with 4 vCPU, 16 GB RAM, Dedicated CPU (AMD)
	CCX33 CloudInstance = "ccx33" // CCX33 is an instance with 8 vCPU, 32 GB RAM, Dedicated CPU (AMD)
	CCX43 CloudInstance = "ccx43" // CCX43 is an instance with 16 vCPU, 64 GB RAM, Dedicated CPU (AMD)
	CCX53 CloudInstance = "ccx53" // CCX53 is an instance with 32 vCPU, 128 GB RAM, Dedicated CPU (AMD)
	CCX63 CloudInstance = "ccx63" // CCX63 is an instance with 48 vCPU, 192 GB RAM, Dedicated CPU (AMD)
)

// Architectures is a struct representing the architectures supported by the cloud provider.
type Architectures struct {
	ARM bool // ARM is a boolean indicating if the ARM architecture is supported
	X86 bool // X86 is a boolean indicating if the X86 architecture is supported
}
