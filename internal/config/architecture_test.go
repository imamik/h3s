package config

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"testing"
)

func TestGetArchitecture(t *testing.T) {
	cases := []struct {
		name   string
		inst   CloudInstance
		expect hcloud.Architecture
	}{
		{"ARM instance", CloudInstance("cax21"), hcloud.ArchitectureARM},
		{"X86 instance", CloudInstance("cx31"), hcloud.ArchitectureX86},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetArchitecture(tc.inst)
			if got != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, got)
			}
		})
	}
}

func TestGetArchitectures(t *testing.T) {
	cfg := &Config{
		ControlPlane: ControlPlane{Pool: NodePool{Instance: CloudInstance("cax21")}},
		WorkerPools:  []NodePool{{Instance: CloudInstance("cx31")}},
	}
	arch := GetArchitectures(cfg)
	if !arch.ARM || !arch.X86 {
		t.Errorf("expected both ARM and X86 to be true, got %+v", arch)
	}

	cfg = &Config{
		ControlPlane: ControlPlane{Pool: NodePool{Instance: CloudInstance("cx31")}},
		WorkerPools:  []NodePool{{Instance: CloudInstance("cx31")}},
	}
	arch = GetArchitectures(cfg)
	if arch.ARM || !arch.X86 {
		t.Errorf("expected only X86 to be true, got %+v", arch)
	}

	cfg = &Config{
		ControlPlane: ControlPlane{Pool: NodePool{Instance: CloudInstance("cax21")}},
		WorkerPools:  []NodePool{{Instance: CloudInstance("cax21")}},
	}
	arch = GetArchitectures(cfg)
	if !arch.ARM || arch.X86 {
		t.Errorf("expected only ARM to be true, got %+v", arch)
	}
}
