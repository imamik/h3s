package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestWorkflowDeletionAdvanced(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	cmd := exec.Command("go", "run", "../../main.go", "destroy", "cluster", "--name=test-cluster")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to destroy cluster: %v\nOutput: %s", err, string(out))
	}
	// Advanced assertion: Ensure cluster is no longer listed
	cmd = exec.Command("go", "run", "../../main.go", "get")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get clusters after deletion: %v\nOutput: %s", err, string(out))
	}
	if strings.Contains(string(out), "test-cluster") {
		t.Errorf("Cluster 'test-cluster' still present after deletion. Output: %s", string(out))
	}
}
