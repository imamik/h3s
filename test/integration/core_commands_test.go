package integration

import (
	"os"
	"os/exec"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	cmd := exec.Command("go", "run", "../../main.go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v, output: %s", err, output)
	}
	if len(output) == 0 || !(contains(string(output), "version")) {
		t.Errorf("expected version output, got: %s", output)
	}
}

func TestInvalidCommand(t *testing.T) {
	if os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1" {
		t.Skip("Skipping real Hetzner integration test (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	cmd := exec.Command("go", "run", "../../main.go", "nonexistentcommand")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for invalid command, got none, output: %s", output)
	}
	if len(output) == 0 || !(contains(string(output), "unknown command") || contains(string(output), "Unknown command")) {
		t.Errorf("expected unknown command error, got: %s", output)
	}
}

// contains checks if substr is in s
// func contains(s, substr string) bool {
// 	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr))))
// }
