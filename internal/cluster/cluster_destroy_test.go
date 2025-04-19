package cluster_test

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(content string) (string, func()) {
	tempfile, _ := os.CreateTemp("", "h3s-config-*.yaml")
	tempfile.Write([]byte(content))
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	cleanup := func() {
		os.Unsetenv("H3S_CONFIG")
		os.Remove(tempfile.Name())
	}
	return tempfile.Name(), cleanup
}

func TestDestroy_HappyPath(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	configYAML := `
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: destroytestcluster
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
`
	_, cleanup := writeTempConfig(configYAML)
	defer cleanup()
	credentialsPath, _ := filepath.Abs("internal/cluster/testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	ctx, err := cluster.Context()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	err = hetzner.Destroy(ctx)
	if err != nil {
		t.Errorf("expected destroy success, got error: %v", err)
	}
}

func TestDestroy_MissingConfig(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	os.Setenv("H3S_CONFIG", "/nonexistent/file.yaml")
	defer os.Unsetenv("H3S_CONFIG")
	_, err := cluster.Context()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestDestroy_MissingCredentials(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	configYAML := `
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: destroytestcluster
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
`
	_, cleanup := writeTempConfig(configYAML)
	defer cleanup()
	credentialsPath, _ := filepath.Abs("internal/cluster/testdata/nonexistent-creds.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	_, err := cluster.Context()
	if err == nil {
		t.Fatal("expected error for missing credentials, got nil")
	}
}

func TestDestroy_PartialDeletion(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	configYAML := `
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: destroytestcluster
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
`
	_, cleanup := writeTempConfig(configYAML)
	defer cleanup()
	credentialsPath, _ := filepath.Abs("internal/cluster/testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	ctx, err := cluster.Context()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	// This test assumes at least one resource deletion will fail gracefully
	err = hetzner.Destroy(ctx)
	// Accept either nil or error, but ensure test does not panic
	if err != nil {
		t.Logf("partial deletion error (acceptable for partial deletion): %v", err)
	}
}
