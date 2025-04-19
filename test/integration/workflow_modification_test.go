package integration

import (
	"os"
	"os/exec"
	"testing"
)

func TestWorkflowModification(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	// Write initial config
	config := []byte("cluster:\n  name: test-cluster\n  region: us-west\n  value: original\n")
	if err := os.WriteFile("h3s.yaml", config, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	defer os.Remove("h3s.yaml")

	cmd := exec.Command("go", "run", "../../main.go", "create", "cluster")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create cluster: %v\nOutput: %s", err, string(out))
	}
	cmd = exec.Command("go", "run", "../../main.go", "get")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get clusters after creation: %v\nOutput: %s", err, string(out))
	}
	if !contains(string(out), "test-cluster") {
		t.Errorf("Cluster creation not reflected. Output: %s", string(out))
	}
	// Modify config
	modConfig := []byte("cluster:\n  name: test-cluster\n  region: us-west\n  value: modified\n")
	if err := os.WriteFile("h3s.yaml", modConfig, 0644); err != nil {
		t.Fatalf("failed to write modified config: %v", err)
	}
	cmd = exec.Command("go", "run", "../../main.go", "create", "cluster")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to modify cluster: %v\nOutput: %s", err, string(out))
	}
	cmd = exec.Command("go", "run", "../../main.go", "get")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get clusters after modification: %v\nOutput: %s", err, string(out))
	}
	if !contains(string(out), "modified") {
		t.Errorf("Cluster modification not reflected. Output: %s", string(out))
	}
}

func TestWorkflowDeletion(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	cmd := exec.Command("h3s", "cluster", "delete", "--name=test-cluster")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to delete cluster: %v\nOutput: %s", err, string(out))
	}
	// TODO: Add assertions to verify cluster deletion
}

func TestWorkflowErrorRecovery(t *testing.T) {
	// Write invalid config
	badConfig := []byte("invalid: yaml\n")
	if err := os.WriteFile("h3s.yaml", badConfig, 0644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}
	defer os.Remove("h3s.yaml")

	cmd := exec.Command("go", "run", "../../main.go", "create", "cluster")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error for invalid cluster creation, got none. Output: %s", string(out))
	}
	// Advanced assertion: Ensure error message is informative
	if !contains(string(out), "Error") && !contains(string(out), "failed") && !contains(string(out), "validation") {
		t.Errorf("Error output not informative: %s", string(out))
	}
}
