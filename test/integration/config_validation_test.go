// Package integration provides integration tests for the h3s application.
package integration

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigValidationIntegration is the main test function for config validation integration tests.
// It sets up the test environment and runs the subtests.
func TestConfigValidationIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "h3s-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run subtests
	t.Run("ValidConfig", func(t *testing.T) {
		testValidConfig(t, tempDir)
	})

	t.Run("MinimalConfig", func(t *testing.T) {
		testMinimalConfig(t, tempDir)
	})

	t.Run("ConfigWithEnvVars", func(t *testing.T) {
		testConfigWithEnvVars(t, tempDir)
	})

	t.Run("MultipleWorkerPools", func(t *testing.T) {
		testMultipleWorkerPools(t, tempDir)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		testMissingRequiredFields(t, tempDir)
	})

	t.Run("InvalidFieldValues", func(t *testing.T) {
		testInvalidFieldValues(t, tempDir)
	})

	t.Run("MalformedYAML", func(t *testing.T) {
		testMalformedYAML(t, tempDir)
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		testNonExistentFile(t, tempDir)
	})

	t.Run("PermissionIssues", func(t *testing.T) {
		testPermissionIssues(t, tempDir)
	})
}

// Helper functions for creating test config files

// createConfigFile creates a config file with the given content in the specified directory.
// Returns the full path to the created file.
func createConfigFile(t *testing.T, dir, filename, content string) string {
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	return path
}

// validConfigYAML returns a valid config YAML string.
func validConfigYAML() string {
	return `ssh_key_paths:
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
}

// minimalConfigYAML returns a minimal valid config YAML string.
func minimalConfigYAML() string {
	return `ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: testcluster
domain: example.com
cert_manager:
  email: test@example.com
  production: false
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`
}

// multipleWorkerPoolsYAML returns a config YAML string with multiple worker pools.
func multipleWorkerPoolsYAML() string {
	return `ssh_key_paths:
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
  - instance: cx41
    location: fsn1
    name: workerpool2
    nodes: 2
  - instance: cx51
    location: hel1
    name: workerpool3
    nodes: 3
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`
}

// missingRequiredFieldsYAML returns a config YAML string with missing required fields.
func missingRequiredFieldsYAML() string {
	return `ssh_key_paths:
  private_key_path: /tmp/id_rsa
  # Missing public_key_path
network_zone: nbg1
k3s_version: v1.28.0
# Missing name
domain: example.com
# Missing cert_manager
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`
}

// invalidFieldValuesYAML returns a config YAML string with invalid field values.
func invalidFieldValuesYAML() string {
	return `ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: invalid-zone
k3s_version: v1.28.0
name: test_cluster # Invalid name (should be kebab-case)
domain: example.com
cert_manager:
  email: invalid-email # Invalid email
  production: false
worker_pools:
  - instance: cx31
    location: nbg1
    name: workerpool1
    nodes: 0 # Invalid node count (should be >= 1)
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`
}

// malformedYAML returns a malformed YAML string.
func malformedYAML() string {
	return `ssh_key_paths:
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
  unclosed: [
`
}

// validCredsYAML returns a valid credentials YAML string.
func validCredsYAML() string {
	return `hcloud_token: p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde
hetzner_dns_token: 1234567890abcdef1234567890abcdef
k3s_token: k3s1234567890abcdef1234567890abcdef
`
}

// Happy path test implementations

// testValidConfig tests loading a valid config file with all required fields.
func testValidConfig(t *testing.T, tempDir string) {
	// Create a valid config file
	configPath := createConfigFile(t, tempDir, "valid-config.yaml", validConfigYAML())

	// Create a credentials file
	credsPath := createConfigFile(t, tempDir, "valid-creds.yaml", validCredsYAML())

	// Set the environment variables to point to the test files
	oldConfigEnv := os.Getenv("H3S_CONFIG")
	oldCredsEnv := os.Getenv("H3S_SECRETS")
	os.Setenv("H3S_CONFIG", configPath)
	os.Setenv("H3S_SECRETS", credsPath)
	defer func() {
		os.Setenv("H3S_CONFIG", oldConfigEnv)
		os.Setenv("H3S_SECRETS", oldCredsEnv)
	}()

	// Run the CLI command to load and validate the config
	out, err := runCLI("version")

	// Check that the command succeeded
	if err != nil {
		t.Errorf("Expected command to succeed with valid config, got error: %v\nOutput: %s", err, out)
	}
}

// testMinimalConfig tests loading a valid config file with minimal required fields.
func testMinimalConfig(t *testing.T, tempDir string) {
	// Create a minimal valid config file
	configPath := createConfigFile(t, tempDir, "minimal-config.yaml", minimalConfigYAML())

	// Create a credentials file
	credsPath := createConfigFile(t, tempDir, "minimal-creds.yaml", validCredsYAML())

	// Set the environment variables to point to the test files
	oldConfigEnv := os.Getenv("H3S_CONFIG")
	oldCredsEnv := os.Getenv("H3S_SECRETS")
	os.Setenv("H3S_CONFIG", configPath)
	os.Setenv("H3S_SECRETS", credsPath)
	defer func() {
		os.Setenv("H3S_CONFIG", oldConfigEnv)
		os.Setenv("H3S_SECRETS", oldCredsEnv)
	}()

	// Run the CLI command to load and validate the config
	out, err := runCLI("version")

	// Check that the command succeeded
	if err != nil {
		t.Errorf("Expected command to succeed with minimal config, got error: %v\nOutput: %s", err, out)
	}
}

// testConfigWithEnvVars tests loading a valid config file with environment variables.
func testConfigWithEnvVars(t *testing.T, tempDir string) {
	// Create a valid config file
	configPath := createConfigFile(t, tempDir, "env-var-config.yaml", validConfigYAML())

	// Create a credentials file
	credsPath := createConfigFile(t, tempDir, "env-var-creds.yaml", validCredsYAML())

	// Set environment variables
	oldConfigEnv := os.Getenv("H3S_CONFIG")
	oldCredsEnv := os.Getenv("H3S_SECRETS")
	oldHetznerEnv := os.Getenv("H3S_HETZNER_ENDPOINT")

	os.Setenv("H3S_CONFIG", configPath)
	os.Setenv("H3S_SECRETS", credsPath)
	// Set a mock Hetzner API endpoint to avoid real API calls
	os.Setenv("H3S_HETZNER_ENDPOINT", "http://localhost:12345")

	defer func() {
		os.Setenv("H3S_CONFIG", oldConfigEnv)
		os.Setenv("H3S_SECRETS", oldCredsEnv)
		os.Setenv("H3S_HETZNER_ENDPOINT", oldHetznerEnv)
	}()

	// Run the CLI command to load and validate the config
	out, err := runCLI("version")

	// Check that the command succeeded
	if err != nil {
		t.Errorf("Expected command to succeed with config and env vars, got error: %v\nOutput: %s", err, out)
	}
}

// testMultipleWorkerPools tests loading a valid config file with different worker pool configurations.
func testMultipleWorkerPools(t *testing.T, tempDir string) {
	// Create a config file with multiple worker pools
	configPath := createConfigFile(t, tempDir, "multi-worker-config.yaml", multipleWorkerPoolsYAML())

	// Create a credentials file
	credsPath := createConfigFile(t, tempDir, "multi-worker-creds.yaml", validCredsYAML())

	// Set the environment variables to point to the test files
	oldConfigEnv := os.Getenv("H3S_CONFIG")
	oldCredsEnv := os.Getenv("H3S_SECRETS")
	os.Setenv("H3S_CONFIG", configPath)
	os.Setenv("H3S_SECRETS", credsPath)
	defer func() {
		os.Setenv("H3S_CONFIG", oldConfigEnv)
		os.Setenv("H3S_SECRETS", oldCredsEnv)
	}()

	// Run the CLI command to load and validate the config
	out, err := runCLI("version")

	// Check that the command succeeded
	if err != nil {
		t.Errorf("Expected command to succeed with multiple worker pools, got error: %v\nOutput: %s", err, out)
	}
}

// Error case test implementations

// testMissingRequiredFields tests loading a config file with missing required fields.
func testMissingRequiredFields(t *testing.T, tempDir string) {
	// Create a config file with missing required fields
	configPath := createConfigFile(t, tempDir, "missing-fields-config.yaml", missingRequiredFieldsYAML())

	// Set the environment variable to point to the test config file
	oldEnv := os.Getenv("H3S_CONFIG")
	os.Setenv("H3S_CONFIG", configPath)
	defer os.Setenv("H3S_CONFIG", oldEnv)

	// Run the CLI command to load and validate the config
	// We expect this to fail with a validation error
	out, err := runCLI("get", "kubeconfig")

	// Check that the command failed with a validation error
	if err == nil {
		t.Errorf("Expected command to fail with missing required fields, but it succeeded\nOutput: %s", out)
	}

	// Check that the error message mentions validation or required fields
	if !contains(out, "validation") && !contains(out, "required") {
		t.Errorf("Expected error message to mention validation or required fields, got: %s", out)
	}
}

// testInvalidFieldValues tests loading a config file with invalid field values.
func testInvalidFieldValues(t *testing.T, tempDir string) {
	// Create a config file with invalid field values
	configPath := createConfigFile(t, tempDir, "invalid-values-config.yaml", invalidFieldValuesYAML())

	// Set the environment variable to point to the test config file
	oldEnv := os.Getenv("H3S_CONFIG")
	os.Setenv("H3S_CONFIG", configPath)
	defer os.Setenv("H3S_CONFIG", oldEnv)

	// Run the CLI command to load and validate the config
	// We expect this to fail with a validation error
	out, err := runCLI("get", "kubeconfig")

	// Check that the command failed with a validation error
	if err == nil {
		t.Errorf("Expected command to fail with invalid field values, but it succeeded\nOutput: %s", out)
	}

	// Check that the error message mentions validation
	if !contains(out, "validation") && !contains(out, "invalid") {
		t.Errorf("Expected error message to mention validation or invalid values, got: %s", out)
	}
}

// testMalformedYAML tests loading a config file with malformed YAML.
func testMalformedYAML(t *testing.T, tempDir string) {
	// Create a config file with malformed YAML
	configPath := createConfigFile(t, tempDir, "malformed-yaml-config.yaml", malformedYAML())

	// Set the environment variable to point to the test config file
	oldEnv := os.Getenv("H3S_CONFIG")
	os.Setenv("H3S_CONFIG", configPath)
	defer os.Setenv("H3S_CONFIG", oldEnv)

	// Run the CLI command to load and validate the config
	// We expect this to fail with a YAML parsing error
	out, err := runCLI("get", "kubeconfig")

	// Check that the command failed with a YAML parsing error
	if err == nil {
		t.Errorf("Expected command to fail with malformed YAML, but it succeeded\nOutput: %s", out)
	}

	// Check that the error message mentions YAML or parsing
	if !contains(out, "yaml") && !contains(out, "YAML") && !contains(out, "parse") {
		t.Errorf("Expected error message to mention YAML or parsing, got: %s", out)
	}
}

// testNonExistentFile tests loading a non-existent config file.
func testNonExistentFile(t *testing.T, tempDir string) {
	// Set the environment variable to point to a non-existent file
	nonExistentPath := filepath.Join(tempDir, "non-existent-config.yaml")
	oldEnv := os.Getenv("H3S_CONFIG")
	os.Setenv("H3S_CONFIG", nonExistentPath)
	defer os.Setenv("H3S_CONFIG", oldEnv)

	// Run the CLI command to load and validate the config
	// We expect this to fail with a file not found error
	out, err := runCLI("get", "kubeconfig")

	// Check that the command failed with a file not found error
	if err == nil {
		t.Errorf("Expected command to fail with non-existent file, but it succeeded\nOutput: %s", out)
	}

	// Check that the error message mentions file not found or no such file
	if !contains(out, "no such file") && !contains(out, "not found") && !contains(out, "cannot find") {
		t.Errorf("Expected error message to mention file not found, got: %s", out)
	}
}

// testPermissionIssues tests loading a config file with permission issues.
func testPermissionIssues(t *testing.T, tempDir string) {
	// Skip on Windows as permission handling is different
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	// Create a config file with restricted permissions
	configPath := createConfigFile(t, tempDir, "no-perm-config.yaml", validConfigYAML())

	// Change permissions to make it unreadable
	err := os.Chmod(configPath, 0000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}
	defer func() {
		if chmodErr := os.Chmod(configPath, 0600); chmodErr != nil {
			t.Logf("Warning: failed to restore file permissions: %v", chmodErr)
		}
	}() // Restore permissions for cleanup

	// Set the environment variable to point to the test config file
	oldEnv := os.Getenv("H3S_CONFIG")
	os.Setenv("H3S_CONFIG", configPath)
	defer os.Setenv("H3S_CONFIG", oldEnv)

	// Run the CLI command to load and validate the config
	// We expect this to fail with a permission error
	out, err := runCLI("get", "kubeconfig")

	// Check that the command failed with a permission error
	if err == nil {
		t.Errorf("Expected command to fail with permission issues, but it succeeded\nOutput: %s", out)
	}

	// Check that the error message mentions permission
	if !contains(out, "permission") && !contains(out, "denied") {
		t.Errorf("Expected error message to mention permission denied, got: %s", out)
	}
}

// Note: Using the existing runCLI and contains functions from the integration package
