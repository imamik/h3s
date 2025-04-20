package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	content := []byte(`
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: testcluster
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
`)
	tempfile, err := os.CreateTemp("", "valid-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	if _, writeErr := tempfile.Write(content); writeErr != nil {
		t.Fatalf("failed to write to temp file: %v", writeErr)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")
	cfg, err := Load()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if cfg == nil {
		t.Errorf("expected config, got nil")
	}
}

func TestLoadConfig_MissingFields(t *testing.T) {
	content := []byte("controlPlane:\n  pool:\nworkerPools:\n  - instance: cx31\n")
	tempfile, err := os.CreateTemp("", "missing-fields-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	if _, writeErr := tempfile.Write(content); writeErr != nil {
		t.Fatalf("failed to write to temp file: %v", writeErr)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")
	_, err = Load()
	if err == nil {
		t.Errorf("expected error for missing fields, got nil")
	}
}

func TestLoadConfig_MalformedYAML(t *testing.T) {
	content := []byte("controlPlane: [unclosed\nworkerPools: - instance: cx31\n")
	tempfile, err := os.CreateTemp("", "malformed-yaml-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	if _, writeErr := tempfile.Write(content); writeErr != nil {
		t.Fatalf("failed to write to temp file: %v", writeErr)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")
	_, err = Load()
	if err == nil {
		t.Errorf("expected error for malformed YAML, got nil")
	}
}

func TestLoadConfig_FileAccessError(t *testing.T) {
	os.Setenv("H3S_CONFIG", "/nonexistent/file.yaml")
	defer os.Unsetenv("H3S_CONFIG")
	_, err := Load()
	if err == nil {
		t.Errorf("expected error for nonexistent file, got nil")
	}
}

func TestLoadConfig_PermissionDenied(t *testing.T) {
	tempfile, err := os.CreateTemp("", "no-perm-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	tempfile.Close()
	if chmodErr := os.Chmod(tempfile.Name(), 0000); chmodErr != nil {
		t.Fatalf("failed to change file permissions: %v", chmodErr)
	}
	defer func() {
		if restoreErr := os.Chmod(tempfile.Name(), 0600); restoreErr != nil {
			t.Logf("Warning: failed to restore file permissions: %v", restoreErr)
		}
	}()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")
	_, err = Load()
	if err == nil {
		t.Errorf("expected error for permission denied, got nil")
	}
}
