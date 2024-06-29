package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

func GetArchitecture(instance CloudInstance) hcloud.Architecture {
	if instance[:3] == "cax" {
		return hcloud.ArchitectureARM
	} else {
		return hcloud.ArchitectureX86
	}
}

func GetArchitectures(config Config) Architectures {
	architectures := Architectures{
		ARM: GetArchitecture(config.ControlPlane.Pool.Instance) == hcloud.ArchitectureARM,
	}
	architectures.X86 = !architectures.ARM

	for _, pool := range config.WorkerPools {
		if GetArchitecture(pool.Instance) == hcloud.ArchitectureARM {
			architectures.ARM = true
		} else {
			architectures.X86 = true
		}
		if architectures.ARM && architectures.X86 {
			break
		}
	}

	return architectures
}
