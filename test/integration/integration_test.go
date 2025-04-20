package integration

import (
	"os"
	"os/exec"
	"testing"

	"h3s/internal/hetzner/mockhetzner"
)

// runCLI runs the CLI command and returns output and error
func runCLI(args ...string) (string, error) {
	//nolint:gosec // G204: Subprocess launched with variable args, assumed safe in test context
	cmd := exec.Command("go", append([]string{"run", "../.."}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// runCLIWithEnv runs the CLI command with environment variables and returns output and error
func runCLIWithEnv(env []string, args ...string) (string, error) {
	//nolint:gosec // G204: Subprocess launched with variable args, assumed safe in test context
	cmd := exec.Command("go", append([]string{"run", "../.."}, args...)...)
	cmd.Env = append(os.Environ(), env...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestHelpCommand(t *testing.T) {
	out, err := runCLI("--help")
	if err != nil {
		t.Fatalf("expected no error, got %v, output: %s", err, out)
	}
	if len(out) == 0 || !(contains(out, "Usage") || contains(out, "usage")) {
		t.Errorf("expected help output, got: %s", out)
	}
}

// TestAdvancedClusterWorkflows tests cluster creation and deletion
//
//nolint:gocyclo // Complexity acceptable for integration test setup
func TestAdvancedClusterWorkflows(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	// Write config/secrets to integration test directory, use relative paths
	config := []byte(`ssh_key_paths:
  private_key_path: id_ed25519
  public_key_path: id_ed25519.pub
network_zone: fsn1
k3s_version: v1.28.2+k3s1
name: test-cluster
domain: test.domain
cert_manager:
  email: test@example.com
  production: false
worker_pools:
  - instance: cx21
    location: fsn1
    name: workerpool1
    nodes: 1
control_plane:
  pool:
    instance: cx21
    location: fsn1
    name: cpool
    nodes: 1
  as_worker_pool: false
`)
	configPath := "h3s.yaml"
	if err := os.WriteFile(configPath, config, 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	defer os.Remove(configPath)

	hcloudToken := "p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde" // 64 chars
	dnsToken := "1234567890abcdef1234567890abcdef"                                    // 32 chars
	k3sToken := "k3s1234567890abcdef1234567890abcdef"                                 // dummy, can be any string
	if len(hcloudToken) != 64 || len(dnsToken) != 32 {
		t.Fatalf("token length error: hcloudToken=%d, dnsToken=%d", len(hcloudToken), len(dnsToken))
	}
	secrets := []byte("hcloud_token: " + hcloudToken + "\nhetzner_dns_token: " + dnsToken + "\nk3s_token: " + k3sToken + "\n")
	// Using a constant for the secrets filename
	//nolint:gosec // G101: Potential hardcoded credentials (filename constant in test)
	const secretsFileName = "h3s-secrets.yaml"
	secretsPath := secretsFileName
	if err := os.WriteFile(secretsPath, secrets, 0600); err != nil {
		t.Fatalf("failed to write secrets: %v", err)
	}
	defer os.Remove(secretsPath)

	out, err := runCLI("create", "cluster")
	if err != nil || !(contains(out, "Cluster created") || contains(out, "created successfully")) {
		t.Errorf("expected cluster creation success, got error: %v, output: %s", err, out)
	}

	out, err = runCLI("destroy", "cluster")
	if err != nil || !(contains(out, "Cluster deleted") || contains(out, "deleted successfully")) {
		t.Errorf("expected cluster deletion success, got error: %v, output: %s", err, out)
	}
}

// TestClusterWorkflow_WithMockHetzner covers basic cluster creation, update, and deletion using a mock Hetzner API.
//
//nolint:gocyclo // Complexity acceptable for integration test setup
func TestClusterWorkflow_WithMockHetzner(t *testing.T) {
	if os.Getenv("H3S_USE_MOCK_HETZNER") != "1" {
		t.Skip("Skipping mock Hetzner integration test (set H3S_USE_MOCK_HETZNER=1 to enable)")
	}

	// Start mock Hetzner API server (success mode)
	mockServer := mockhetzner.NewHetznerMockScenario("/v1", "success")
	defer mockServer.Close()

	// Write config/secrets to integration test directory, use relative paths
	config := []byte(`ssh_key_paths:
  private_key_path: id_ed25519
  public_key_path: id_ed25519.pub
network_zone: fsn1
k3s_version: v1.28.2+k3s1
name: test-cluster
domain: test.domain
cert_manager:
  email: test@example.com
  production: false
worker_pools:
  - instance: cx21
    location: fsn1
    name: workerpool1
    nodes: 1
control_plane:
  pool:
    instance: cx21
    location: fsn1
    name: cpool
    nodes: 1
  as_worker_pool: false
hetzner_api_endpoint: "` + mockServer.Server.URL + `"
`)
	configPath := "h3s.yaml"
	if err := os.WriteFile(configPath, config, 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	defer os.Remove(configPath)

	hcloudToken := "p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde" // 64 chars
	dnsToken := "1234567890abcdef1234567890abcdef"                                    // 32 chars
	k3sToken := "k3s1234567890abcdef1234567890abcdef"                                 // dummy, can be any string
	if len(hcloudToken) != 64 || len(dnsToken) != 32 {
		t.Fatalf("token length error: hcloudToken=%d, dnsToken=%d", len(hcloudToken), len(dnsToken))
	}
	secrets := []byte("hcloud_token: " + hcloudToken + "\nhetzner_dns_token: " + dnsToken + "\nk3s_token: " + k3sToken + "\n")
	// Using the same constant as above
	//nolint:gosec // G101: Potential hardcoded credentials (filename constant in test)
	const secretsFileName = "h3s-secrets.yaml"
	secretsPath := secretsFileName
	if err := os.WriteFile(secretsPath, secrets, 0600); err != nil {
		t.Fatalf("failed to write secrets: %v", err)
	}
	defer os.Remove(secretsPath)

	env := []string{"H3S_HETZNER_ENDPOINT=" + mockServer.Server.URL}

	out, err := runCLIWithEnv(env, "create", "cluster")
	if err != nil || !(contains(out, "Cluster created") || contains(out, "created successfully")) {
		t.Errorf("expected cluster creation success, got error: %v, output: %s", err, out)
	}

	out, err = runCLIWithEnv(env, "destroy", "cluster")
	if err != nil || !(contains(out, "Cluster deleted") || contains(out, "deleted successfully")) {
		t.Errorf("expected cluster deletion success, got error: %v, output: %s", err, out)
	}

	// Now test error scenario
	mockServerError := mockhetzner.NewHetznerMockScenario("/v1", "error")
	defer mockServerError.Close()
	configError := []byte(`ssh_key_paths:
  private_key_path: id_ed25519
  public_key_path: id_ed25519.pub
network_zone: fsn1
k3s_version: v1.28.2+k3s1
name: test-cluster-err
domain: test.domain
cert_manager:
  email: test@example.com
  production: false
worker_pools:
  - instance: cx21
    location: fsn1
    name: workerpool1
    nodes: 1
control_plane:
  pool:
    instance: cx21
    location: fsn1
    name: cpool
    nodes: 1
  as_worker_pool: false
hetzner_api_endpoint: "` + mockServerError.Server.URL + `"
`)
	configErrorPath := "h3s-err.yaml"
	if writeErr := os.WriteFile(configErrorPath, configError, 0600); writeErr != nil {
		t.Fatalf("failed to write error config: %v", writeErr)
	}
	defer os.Remove(configErrorPath)
	//nolint:gosec // G101: Potential hardcoded credentials (filename constant in test)
	secretsErrorPath := "h3s-secrets-err.yaml"
	if writeErr := os.WriteFile(secretsErrorPath, secrets, 0600); writeErr != nil {
		t.Fatalf("failed to write error secrets: %v", writeErr)
	}
	defer os.Remove(secretsErrorPath)
	envError := []string{"H3S_HETZNER_ENDPOINT=" + mockServerError.Server.URL, "H3S_CONFIG=" + configErrorPath, "H3S_SECRETS=" + secretsErrorPath}
	out, err = runCLIWithEnv(envError, "create", "cluster")
	if err == nil || !(contains(out, "internal error") || contains(out, "failed") || contains(out, "error")) {
		t.Errorf("expected cluster creation error, got: %v, output: %s", err, out)
	}
}
