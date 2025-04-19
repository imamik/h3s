package components

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"net"
	"testing"
)

func TestVars_ResourceLimits_NotPresent(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/16")
	conf := &config.Config{Domain: "example.com"}
	creds := &credentials.ProjectCredentials{}
	network := &hcloud.Network{Name: "testnet", IPRange: ipnet}
	lb := &hcloud.LoadBalancer{Name: "lb1", Location: &hcloud.Location{Name: "fsn1"}}

	vars := GetVars(conf, creds, network, lb)
	if _, ok := vars["ResourceLimits"]; ok {
		t.Errorf("ResourceLimits key should not be present in vars map yet")
	}
	if _, ok := vars["Quotas"]; ok {
		t.Errorf("Quotas key should not be present in vars map yet")
	}
	// TODO: When resource limits/quotas are implemented, update this test to validate their values
}

func TestVars_UnsupportedOptions(t *testing.T) {
	// TODO: Test that unsupported options are detected and handled
	t.Skip("Unsupported options test scaffold")
}
