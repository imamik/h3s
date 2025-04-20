// Package integration provides end-to-end tests for the h3s CLI.
package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"h3s/internal/hetzner/mockhetzner"
)

// TestE2EWorkflows is the main test function for end-to-end CLI workflow tests.
// It sets up the test environment and runs the subtests.
func TestE2EWorkflows(t *testing.T) {
	// Skip if not explicitly enabled
	if os.Getenv("H3S_ENABLE_E2E_TESTS") != "1" {
		t.Skip("Skipping E2E tests (set H3S_ENABLE_E2E_TESTS=1 to enable)")
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "h3s-e2e-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create SSH key files for testing
	privateKeyPath, publicKeyPath, err := createTestSSHKeys(tempDir)
	if err != nil {
		t.Fatalf("Failed to create test SSH keys: %v", err)
	}

	// Start mock Hetzner API server
	mockServer := mockhetzner.NewHetznerMockScenario("/v1", "success")
	defer mockServer.Close()

	// Set up test environment
	testEnv := &e2eTestEnv{
		tempDir:        tempDir,
		mockServer:     mockServer,
		privateKeyPath: privateKeyPath,
		publicKeyPath:  publicKeyPath,
	}

	// Run subtests
	t.Run("BasicCommands", func(t *testing.T) {
		testBasicCommands(t, testEnv)
	})

	t.Run("ConfigCommands", func(t *testing.T) {
		testConfigCommands(t, testEnv)
	})

	t.Run("ClusterOperations", func(t *testing.T) {
		testClusterOperations(t, testEnv)
	})

	t.Run("CompleteWorkflow", func(t *testing.T) {
		testCompleteWorkflow(t, testEnv)
	})

	t.Run("ErrorRecovery", func(t *testing.T) {
		testErrorRecovery(t, testEnv)
	})
}

// e2eTestEnv holds the test environment for E2E tests
type e2eTestEnv struct {
	tempDir        string
	mockServer     *mockhetzner.MockServer
	privateKeyPath string
	publicKeyPath  string
}

// createTestSSHKeys creates temporary SSH key files for testing
func createTestSSHKeys(dir string) (privateKeyPath, publicKeyPath string, err error) {
	privateKeyPath = filepath.Join(dir, "id_ed25519")
	publicKeyPath = filepath.Join(dir, "id_ed25519.pub")

	// Sample SSH key content for testing
	privateKey := "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\nQyNTUxOQAAACDFIzSM+Yg/xFbZGX7KwJVBK0GQgwQT+xVgWsK1KvYZMAAAAJgUKTd9FCk3\nfQAAAAtzc2gtZWQyNTUxOQAAACDFIzSM+Yg/xFbZGX7KwJVBK0GQgwQT+xVgWsK1KvYZMA\nAAAECTVAYb8PxOIgKZNxTkZwTTQYrMnQxK9CjWbGG8jFGJxsUjNIz5iD/EVtkZfsrAlUEr\nQZCDBBP7FWBawrUq9hkwAAAAEHRlc3RAZXhhbXBsZS5jb20BAgMEBQ==\n-----END OPENSSH PRIVATE KEY-----"
	publicKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMUjNIz5iD/EVtkZfsrAlUErQZCDBBP7FWBawrUq9hkw test@example.com"

	if err := os.WriteFile(privateKeyPath, []byte(privateKey), 0600); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(publicKeyPath, []byte(publicKey), 0600); err != nil {
		return "", "", err
	}

	return privateKeyPath, publicKeyPath, nil
}

// createTestConfig creates a test configuration file
func createTestConfig(dir string, mockServerURL string, privateKeyPath, publicKeyPath string) (string, error) {
	configPath := filepath.Join(dir, "h3s.yaml")
	config := []byte(`ssh_key_paths:
  private_key_path: ` + privateKeyPath + `
  public_key_path: ` + publicKeyPath + `
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
hetzner_api_endpoint: "` + mockServerURL + `"
`)

	return configPath, os.WriteFile(configPath, config, 0600)
}

// Test token constants
const (
	testHCloudToken = "p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde" // 64 chars
	testDNSToken    = "1234567890abcdef1234567890abcdef"                                 // 32 chars
	testK3sToken    = "k3s1234567890abcdef1234567890abcdef"                              // dummy, can be any string
)

// createTestCredentials creates a test credentials file
func createTestCredentials(dir string) (string, error) {
	credsPath := filepath.Join(dir, "h3s-secrets.yaml")
	hcloudToken := testHCloudToken
	dnsToken := testDNSToken
	k3sToken := testK3sToken

	creds := []byte("hcloud_token: " + hcloudToken + "\nhetzner_dns_token: " + dnsToken + "\nk3s_token: " + k3sToken + "\n")

	return credsPath, os.WriteFile(credsPath, creds, 0600)
}

// findProjectRoot finds the root directory of the project
func findProjectRoot() string {
	// Start from the current directory and go up until we find a go.mod file
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory
			break
		}
		dir = parent
	}

	return ""
}

// runCLIWithEnvAndPath runs the CLI command with the given binary path, environment variables, and arguments
func runCLIWithEnvAndPath(binaryPath string, env []string, args ...string) (string, error) {
	// Create the command
	cmd := exec.Command(binaryPath, args...)

	// Set environment variables
	cmd.Env = append(os.Environ(), env...)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Combine stdout and stderr for the output
	output := stdout.String() + stderr.String()

	return output, err
}

// runCLIWithEnvAndDir runs the CLI command with environment variables and working directory
func runCLIWithEnvAndDir(dir string, env []string, args ...string) (string, error) {
	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Change to test directory
	if chErr := os.Chdir(dir); chErr != nil {
		return "", chErr
	}

	// Find the h3s binary
	// First, try the project root directory
	projectRoot := findProjectRoot()
	binaryPath := filepath.Join(projectRoot, "h3s")
	if _, statErr := os.Stat(binaryPath); os.IsNotExist(statErr) {
		// If not found, try the current directory
		binaryPath = filepath.Join(currentDir, "h3s")
		if _, statErr := os.Stat(binaryPath); os.IsNotExist(statErr) {
			// If still not found, use the command name and hope it's in PATH
			binaryPath = "h3s"
		}
	}

	// Set up the command with the binary path
	oldArgs := args
	args = make([]string, len(oldArgs))
	copy(args, oldArgs)

	// Run command with the binary path
	out, err := runCLIWithEnvAndPath(binaryPath, env, args...)

	// Restore original directory
	if chErr := os.Chdir(currentDir); chErr != nil {
		return out, chErr
	}

	return out, err
}

// Basic command tests

// testBasicCommands tests basic CLI commands that don't require configuration
func testBasicCommands(t *testing.T, _ *e2eTestEnv) {
	// Test version command
	t.Run("VersionCommand", func(t *testing.T) {
		out, err := runCLI("version")
		if err != nil {
			t.Fatalf("Expected no error for version command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "version") {
			t.Errorf("Expected version output to contain 'version', got: %s", out)
		}
	})

	// Test help command
	t.Run("HelpCommand", func(t *testing.T) {
		out, err := runCLI("--help")
		if err != nil {
			t.Fatalf("Expected no error for help command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "Usage") {
			t.Errorf("Expected help output to contain 'Usage', got: %s", out)
		}
	})

	// Test invalid command
	t.Run("InvalidCommand", func(t *testing.T) {
		_, err := runCLI("nonexistentcommand")
		if err == nil {
			t.Errorf("Expected error for invalid command, got none")
		}
	})
}

// Configuration command tests

// testConfigCommands tests commands related to configuration management
func testConfigCommands(t *testing.T, env *e2eTestEnv) {
	// Test create config command
	t.Run("CreateConfig", func(t *testing.T) {
		// Run the create config command with interactive mode disabled
		envVars := []string{"H3S_NON_INTERACTIVE=1"}
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "create", "config",
			"--name=test-cluster",
			"--domain=test.domain",
			"--network-zone=fsn1",
			"--k3s-version=v1.28.2+k3s1",
			"--control-plane-nodes=1",
			"--control-plane-instance=cx21",
			"--control-plane-location=fsn1",
			"--worker-nodes=1",
			"--worker-instance=cx21",
			"--worker-location=fsn1",
			"--ssh-private-key-path="+env.privateKeyPath,
			"--ssh-public-key-path="+env.publicKeyPath,
			"--cert-manager-email=test@example.com",
			"--cert-manager-production=false",
			"--hetzner-api-endpoint="+env.mockServer.Server.URL,
		)
		if err != nil {
			t.Fatalf("Expected no error for create config command, got: %v\nOutput: %s", err, out)
		}

		// Check if config file was created
		configPath := filepath.Join(env.tempDir, "h3s.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("Config file was not created at %s", configPath)
		}
	})

	// Test create credentials command
	t.Run("CreateCredentials", func(t *testing.T) {
		// Run the create credentials command with interactive mode disabled
		envVars := []string{"H3S_NON_INTERACTIVE=1"}
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "create", "credentials",
			"--hcloud-token=p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
			"--hetzner-dns-token=1234567890abcdef1234567890abcdef",
			"--k3s-token=k3s1234567890abcdef1234567890abcdef",
		)
		if err != nil {
			t.Fatalf("Expected no error for create credentials command, got: %v\nOutput: %s", err, out)
		}

		// Check if credentials file was created
		credsPath := filepath.Join(env.tempDir, "h3s-secrets.yaml")
		if _, err := os.Stat(credsPath); os.IsNotExist(err) {
			t.Errorf("Credentials file was not created at %s", credsPath)
		}
	})
}

// Cluster operation tests

// testClusterOperations tests individual cluster operations with mocked API
func testClusterOperations(t *testing.T, env *e2eTestEnv) {
	// Create config and credentials files for testing
	configPath, err := createTestConfig(env.tempDir, env.mockServer.Server.URL, env.privateKeyPath, env.publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	credsPath, err := createTestCredentials(env.tempDir)
	if err != nil {
		t.Fatalf("Failed to create test credentials: %v", err)
	}

	// Set environment variables for tests
	envVars := []string{
		"H3S_CONFIG=" + configPath,
		"H3S_SECRETS=" + credsPath,
		"H3S_HETZNER_ENDPOINT=" + env.mockServer.Server.URL,
	}

	// Test create cluster command
	t.Run("CreateCluster", func(t *testing.T) {
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "create", "cluster")
		if err != nil {
			t.Fatalf("Expected no error for create cluster command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "created") && !contains(out, "success") {
			t.Errorf("Expected output to indicate successful creation, got: %s", out)
		}
	})

	// Test get kubeconfig command
	t.Run("GetKubeconfig", func(t *testing.T) {
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "get", "kubeconfig")
		if err != nil {
			t.Fatalf("Expected no error for get kubeconfig command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "apiVersion") && !contains(out, "kubeconfig") {
			t.Errorf("Expected output to contain kubeconfig content, got: %s", out)
		}
	})

	// Test get token command
	t.Run("GetToken", func(t *testing.T) {
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "get", "token")
		if err != nil {
			t.Fatalf("Expected no error for get token command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "token") {
			t.Errorf("Expected output to contain token information, got: %s", out)
		}
	})

	// Test ssh command
	t.Run("SSHCommand", func(t *testing.T) {
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "ssh", "ls", "-la")
		// SSH command might fail in mock environment, but we should at least check it tries to execute
		if err != nil {
			if !contains(out, "ssh") && !contains(out, "command") {
				t.Errorf("Expected output to mention SSH command execution, got: %s", out)
			}
		}
	})

	// Test destroy cluster command
	t.Run("DestroyCluster", func(t *testing.T) {
		out, err := runCLIWithEnvAndDir(env.tempDir, envVars, "destroy", "cluster")
		if err != nil {
			t.Fatalf("Expected no error for destroy cluster command, got: %v\nOutput: %s", err, out)
		}
		if !contains(out, "deleted") && !contains(out, "destroyed") && !contains(out, "success") {
			t.Errorf("Expected output to indicate successful deletion, got: %s", out)
		}
	})
}

// Complete workflow tests

// testCompleteWorkflow tests a complete cluster lifecycle workflow
func testCompleteWorkflow(t *testing.T, env *e2eTestEnv) {
	// Create a clean directory for this test
	workflowDir, err := os.MkdirTemp(env.tempDir, "workflow-*")
	if err != nil {
		t.Fatalf("Failed to create workflow test directory: %v", err)
	}
	defer os.RemoveAll(workflowDir)

	// Set environment variables
	envVars := []string{
		"H3S_NON_INTERACTIVE=1",
		"H3S_HETZNER_ENDPOINT=" + env.mockServer.Server.URL,
	}

	// Step 1: Create configuration
	t.Log("Step 1: Creating configuration")
	out, err := runCLIWithEnvAndDir(workflowDir, envVars, "create", "config",
		"--name=workflow-cluster",
		"--domain=workflow.test",
		"--network-zone=fsn1",
		"--k3s-version=v1.28.2+k3s1",
		"--control-plane-nodes=1",
		"--control-plane-instance=cx21",
		"--control-plane-location=fsn1",
		"--worker-nodes=1",
		"--worker-instance=cx21",
		"--worker-location=fsn1",
		"--ssh-private-key-path="+env.privateKeyPath,
		"--ssh-public-key-path="+env.publicKeyPath,
		"--cert-manager-email=workflow@example.com",
		"--cert-manager-production=false",
		"--hetzner-api-endpoint="+env.mockServer.Server.URL,
	)
	if err != nil {
		t.Fatalf("Step 1 failed: %v\nOutput: %s", err, out)
	}

	// Step 2: Create credentials
	t.Log("Step 2: Creating credentials")
	out, err = runCLIWithEnvAndDir(workflowDir, envVars, "create", "credentials",
		"--hcloud-token=p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
		"--hetzner-dns-token=1234567890abcdef1234567890abcdef",
		"--k3s-token=k3s1234567890abcdef1234567890abcdef",
	)
	if err != nil {
		t.Fatalf("Step 2 failed: %v\nOutput: %s", err, out)
	}

	// Update environment variables to include config and secrets paths
	configPath := filepath.Join(workflowDir, "h3s.yaml")
	credsPath := filepath.Join(workflowDir, "h3s-secrets.yaml")
	envVars = append(envVars,
		"H3S_CONFIG="+configPath,
		"H3S_SECRETS="+credsPath,
	)

	// Step 3: Create cluster
	t.Log("Step 3: Creating cluster")
	out, err = runCLIWithEnvAndDir(workflowDir, envVars, "create", "cluster")
	if err != nil {
		t.Fatalf("Step 3 failed: %v\nOutput: %s", err, out)
	}

	// Step 4: Get kubeconfig
	t.Log("Step 4: Getting kubeconfig")
	out, err = runCLIWithEnvAndDir(workflowDir, envVars, "get", "kubeconfig")
	if err != nil {
		t.Fatalf("Step 4 failed: %v\nOutput: %s", err, out)
	}

	// Step 5: Execute SSH command
	t.Log("Step 5: Executing SSH command")
	_, _ = runCLIWithEnvAndDir(workflowDir, envVars, "ssh", "ls", "-la")
	// SSH might fail in mock environment, so we don't check for errors

	// Step 6: Destroy cluster
	t.Log("Step 6: Destroying cluster")
	out, err = runCLIWithEnvAndDir(workflowDir, envVars, "destroy", "cluster")
	if err != nil {
		t.Fatalf("Step 6 failed: %v\nOutput: %s", err, out)
	}

	t.Log("Complete workflow test succeeded")
}

// Error recovery tests

// testErrorRecovery tests recovery from error conditions
func testErrorRecovery(t *testing.T, env *e2eTestEnv) {
	// Create a clean directory for this test
	recoveryDir, err := os.MkdirTemp(env.tempDir, "recovery-*")
	if err != nil {
		t.Fatalf("Failed to create recovery test directory: %v", err)
	}
	defer os.RemoveAll(recoveryDir)

	// Test recovery from invalid configuration
	t.Run("RecoveryFromInvalidConfig", func(t *testing.T) {
		// Create an invalid config file
		invalidConfigPath := filepath.Join(recoveryDir, "invalid-config.yaml")
		invalidConfig := []byte(`invalid: yaml: content
this is not valid yaml`)
		if err := os.WriteFile(invalidConfigPath, invalidConfig, 0600); err != nil {
			t.Fatalf("Failed to create invalid config file: %v", err)
		}

		// Create valid credentials
		credsPath, err := createTestCredentials(recoveryDir)
		if err != nil {
			t.Fatalf("Failed to create test credentials: %v", err)
		}

		// Try to create cluster with invalid config
		envVars := []string{
			"H3S_CONFIG=" + invalidConfigPath,
			"H3S_SECRETS=" + credsPath,
		}
		out, err := runCLIWithEnvAndDir(recoveryDir, envVars, "create", "cluster")

		// Should fail
		if err == nil {
			t.Errorf("Expected error with invalid config, got success\nOutput: %s", out)
		}

		// Now create a valid config file
		validConfigPath, err := createTestConfig(recoveryDir, env.mockServer.Server.URL, env.privateKeyPath, env.publicKeyPath)
		if err != nil {
			t.Fatalf("Failed to create valid config: %v", err)
		}

		// Try again with valid config
		envVars = []string{
			"H3S_CONFIG=" + validConfigPath,
			"H3S_SECRETS=" + credsPath,
			"H3S_HETZNER_ENDPOINT=" + env.mockServer.Server.URL,
		}
		out, err = runCLIWithEnvAndDir(recoveryDir, envVars, "create", "cluster")

		// Should succeed
		if err != nil {
			t.Errorf("Expected success after fixing config, got error: %v\nOutput: %s", err, out)
		}

		// Clean up
		out, err = runCLIWithEnvAndDir(recoveryDir, envVars, "destroy", "cluster")
		if err != nil {
			t.Logf("Warning: Cleanup failed: %v\nOutput: %s", err, out)
		}
	})

	// Test recovery from API errors
	t.Run("RecoveryFromAPIError", func(t *testing.T) {
		// Create config and credentials
		configPath, err := createTestConfig(recoveryDir, "http://invalid-url", env.privateKeyPath, env.publicKeyPath)
		if err != nil {
			t.Fatalf("Failed to create test config: %v", err)
		}
		credsPath, err := createTestCredentials(recoveryDir)
		if err != nil {
			t.Fatalf("Failed to create test credentials: %v", err)
		}

		// Try to create cluster with invalid API endpoint
		envVars := []string{
			"H3S_CONFIG=" + configPath,
			"H3S_SECRETS=" + credsPath,
			"H3S_HETZNER_ENDPOINT=http://invalid-url",
		}
		out, err := runCLIWithEnvAndDir(recoveryDir, envVars, "create", "cluster")

		// Should fail
		if err == nil {
			t.Errorf("Expected error with invalid API endpoint, got success\nOutput: %s", out)
		}

		// Now update the config with valid API endpoint
		envVars = []string{
			"H3S_CONFIG=" + configPath,
			"H3S_SECRETS=" + credsPath,
			"H3S_HETZNER_ENDPOINT=" + env.mockServer.Server.URL,
		}
		out, err = runCLIWithEnvAndDir(recoveryDir, envVars, "create", "cluster")

		// Should succeed
		if err != nil {
			t.Errorf("Expected success after fixing API endpoint, got error: %v\nOutput: %s", err, out)
		}

		// Clean up
		out, err = runCLIWithEnvAndDir(recoveryDir, envVars, "destroy", "cluster")
		if err != nil {
			t.Logf("Warning: Cleanup failed: %v\nOutput: %s", err, out)
		}
	})

	// Test configuration modification workflow
	t.Run("ConfigModification", func(t *testing.T) {
		// Create a subdirectory for this test
		configModDir := filepath.Join(recoveryDir, "config-mod")
		if err := os.MkdirAll(configModDir, 0755); err != nil {
			t.Fatalf("Failed to create config modification test directory: %v", err)
		}

		// Create initial config and credentials
		configPath, err := createTestConfig(configModDir, env.mockServer.Server.URL, env.privateKeyPath, env.publicKeyPath)
		if err != nil {
			t.Fatalf("Failed to create test config: %v", err)
		}
		credsPath, err := createTestCredentials(configModDir)
		if err != nil {
			t.Fatalf("Failed to create test credentials: %v", err)
		}

		// Set environment variables
		envVars := []string{
			"H3S_CONFIG=" + configPath,
			"H3S_SECRETS=" + credsPath,
			"H3S_HETZNER_ENDPOINT=" + env.mockServer.Server.URL,
		}

		// Create cluster with initial config
		out, err := runCLIWithEnvAndDir(configModDir, envVars, "create", "cluster")
		if err != nil {
			t.Fatalf("Failed to create cluster with initial config: %v\nOutput: %s", err, out)
		}

		// Modify the configuration (add a worker pool)
		// Read the current config
		configData, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Add a new worker pool to the config
		modifiedConfig := string(configData) + `
worker_pools:
  - instance: cx31
    location: hel1
    name: new-workers
    nodes: 2
`

		// Write the modified config back
		if writeErr := os.WriteFile(configPath, []byte(modifiedConfig), 0600); writeErr != nil {
			t.Fatalf("Failed to write modified config: %v", writeErr)
		}

		// Update the cluster with the modified config
		// In a real implementation, there would be an "update cluster" command
		// For now, we'll just destroy and recreate
		out, err = runCLIWithEnvAndDir(configModDir, envVars, "destroy", "cluster")
		if err != nil {
			t.Fatalf("Failed to destroy cluster for update: %v\nOutput: %s", err, out)
		}

		out, err = runCLIWithEnvAndDir(configModDir, envVars, "create", "cluster")
		if err != nil {
			t.Fatalf("Failed to recreate cluster with modified config: %v\nOutput: %s", err, out)
		}

		// Clean up
		out, err = runCLIWithEnvAndDir(configModDir, envVars, "destroy", "cluster")
		if err != nil {
			t.Logf("Warning: Cleanup failed: %v\nOutput: %s", err, out)
		}
	})
}
