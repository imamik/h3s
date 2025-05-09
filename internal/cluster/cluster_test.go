package cluster

import (
	"fmt"
	"h3s/internal/config"
	"os"
	"path/filepath"
	"testing"
)

// testConfigYAML is a common test configuration used across tests
const testConfigYAML = `
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
`

func writeTempConfig(content string) (string, func()) {
	tempfile, _ := os.CreateTemp("", "h3s-config-*.yaml")
	_, err := tempfile.Write([]byte(content))
	if err != nil {
		panic(err)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	cleanup := func() {
		os.Unsetenv("H3S_CONFIG")
		os.Remove(tempfile.Name())
	}
	return tempfile.Name(), cleanup
}

func TestContext_Success(t *testing.T) {
	credentialsPath, _ := filepath.Abs("testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	_, cleanup := writeTempConfig(testConfigYAML)
	defer cleanup()
	_, err := Context()
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestContext_MissingConfig(t *testing.T) {
	os.Setenv("H3S_CONFIG", "/nonexistent/file.yaml")
	defer os.Unsetenv("H3S_CONFIG")
	_, err := Context()
	if err == nil {
		t.Error("expected error for missing config, got nil")
	}
}

func TestContext_MissingCredentials(t *testing.T) {
	_, cleanup := writeTempConfig(testConfigYAML)
	defer cleanup()
	// Temporarily rename secrets file if it exists
	tmpDir, err := os.MkdirTemp("", "test-context-missing-credentials")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	// Create necessary files
	configPath := filepath.Join(tmpDir, "h3s.yaml")
	_ = configPath // Assign to blank identifier to silence unused warning
	// Set credentials path to a truly non-existent file
	credentialsPath := filepath.Join(tmpDir, "definitely-does-not-exist.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	_, err = Context()
	if err == nil {
		t.Error("expected error for missing credentials, got nil")
	}
}

func TestContext_Idempotency(t *testing.T) {
	credentialsPath, _ := filepath.Abs("testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	_, cleanup := writeTempConfig(testConfigYAML)
	defer cleanup()
	for i := 0; i < 5; i++ {
		_, err := Context()
		if err != nil {
			t.Errorf("idempotency failed on iteration %d: %v", i, err)
		}
	}
}

func TestContext_LargeCluster(t *testing.T) {
	credentialsPath, _ := filepath.Abs("testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	configYAML := `
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
` + generateLargeWorkerPools(100) + `
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 3
  as_worker_pool: false
`
	_, cleanup := writeTempConfig(configYAML)
	defer cleanup()
	_, err := Context()
	if err != nil {
		t.Errorf("large cluster config failed: %v", err)
	}
}

func generateLargeWorkerPools(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += fmt.Sprintf("  - instance: cx31\n    location: nbg1\n    name: workerpool%03d\n    nodes: 1\n", i)
	}
	return result
}

// TestContext_ResourceExhaustion simulates too many nodes (resource exhaustion).
func TestContext_ResourceExhaustion(t *testing.T) {
	credentialsPath, _ := filepath.Abs("testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)
	configYAML := `
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
    nodes: 10000
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
	_, err := Context()
	// Note: resource exhaustion is not enforced by the current implementation, so we expect success.
	if err != nil {
		t.Errorf("expected success (resource exhaustion not enforced), got error: %v", err)
	}
}

func TestContext_Recovery(t *testing.T) {
	// Start with non-existent paths
	os.Setenv("H3S_CONFIG", "/nonexistent/file.yaml")
	os.Setenv("H3S_CREDENTIALS", "/nonexistent/creds.yaml")
	_, err := Context()
	if err == nil {
		t.Error("expected error for missing config/creds, got nil")
	}
	configPath, cleanup := writeTempConfig(testConfigYAML)
	defer cleanup()
	os.Setenv("H3S_CONFIG", configPath)
	// Correctly resolve and set the valid credentials path *after* fixing the config
	credentialsPath, _ := filepath.Abs("testdata/valid-credentials.yaml")
	os.Setenv("H3S_CREDENTIALS", credentialsPath)

	_, err = Context()
	if err != nil {
		t.Errorf("expected recovery after fixing config/creds, got error: %v", err)
	}
}

func TestWriteTempConfig_Cleanup(t *testing.T) {
	configYAML := "foo: bar"
	path, cleanup := writeTempConfig(configYAML)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("temp config not created: %v", err)
	}
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("temp config not cleaned up: %v", err)
	}
}

func TestGetLabels_DefaultAndAdditional(t *testing.T) {
	c := &Cluster{Config: &config.Config{Name: "testproject"}}
	labels := c.GetLabels()
	if labels["project"] != "testproject" || labels["origin"] != "h3s" {
		t.Errorf("default labels incorrect: %+v", labels)
	}
	add := map[string]string{"env": "dev", "tier": "backend"}
	labels = c.GetLabels(add)
	if labels["env"] != "dev" || labels["tier"] != "backend" {
		t.Errorf("additional labels not set: %+v", labels)
	}
}

func TestGetName_Combinations(t *testing.T) {
	c := &Cluster{Config: &config.Config{Name: "myproj"}}
	name := c.GetName()
	if name != "myproj-h3s" {
		t.Errorf("expected 'myproj-h3s', got '%s'", name)
	}
	name = c.GetName("db", "prod")
	if name != "myproj-h3s-db-prod" {
		t.Errorf("expected 'myproj-h3s-db-prod', got '%s'", name)
	}
}

func TestPrintWorkingDirectory(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	t.Logf("Current working directory: %s", dir)
}

// Additional edge case and idempotency tests can be added as needed.
