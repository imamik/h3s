package config

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
	Shared    CloudInstanceType = "shared"
	Dedicated CloudInstanceType = "dedicated"
)

type CloudInstanceArchitecture string

const (
	Intel CloudInstanceArchitecture = "intel"
	AMD   CloudInstanceArchitecture = "amd"
	Arm   CloudInstanceArchitecture = "arm"
)

type CloudInstance string

// Shared vCPU AMD
const (
	CPX11 CloudInstance = "CPX11"
	CPX21 CloudInstance = "CPX21"
	CPX31 CloudInstance = "CPX31"
	CPX41 CloudInstance = "CPX41"
	CPX51 CloudInstance = "CPX51"
)

// Shared vCPU Intel
const (
	CX22 CloudInstance = "CX22"
	CX32 CloudInstance = "CX32"
	CX42 CloudInstance = "CX42"
	CX52 CloudInstance = "CX52"
)

// Shared vCPU Arm
const (
	CAX11 CloudInstance = "CAX11"
	CAX21 CloudInstance = "CAX21"
	CAX31 CloudInstance = "CAX31"
	CAX41 CloudInstance = "CAX41"
)

// Dedicated vCPU AMD
const (
	CCX13 CloudInstance = "CCX13"
	CCX23 CloudInstance = "CCX23"
	CCX33 CloudInstance = "CCX33"
	CCX43 CloudInstance = "CCX43"
	CCX53 CloudInstance = "CCX53"
	CCX63 CloudInstance = "CCX63"
)
