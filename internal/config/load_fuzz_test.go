//go:build go1.18
// +build go1.18

package config

import (
	"os"
	"testing"
)

func FuzzLoadConfig(f *testing.F) {
	// Add a valid config seed
	f.Add([]byte(`ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: fuzzcluster
domain: example.com
cert_manager:
  email: test@example.com
  production: false
worker_pools:
  - instance: cx31
    location: nbg1
    name: workerpool1
    nodes: 1
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Write fuzz data to a temp file
		tmp, err := os.CreateTemp("", "h3s-fuzz-*.yaml")
		if err != nil {
			t.Skip("could not create temp file")
		}
		defer os.Remove(tmp.Name())
		_, _ = tmp.Write(data)
		tmp.Close()
		os.Setenv("H3S_CONFIG", tmp.Name())
		defer os.Unsetenv("H3S_CONFIG")
		_, _ = Load() // We only care about panics, not errors
	})
}
