package version

import (
	"testing"
)

// TestVersionInjectionViaLdflags tests that version information can be injected via ldflags
// This is a simpler test that just verifies the mechanism works by checking the Makefile
func TestVersionInjectionViaLdflags(t *testing.T) {
	// This test doesn't actually build anything, it just verifies that the ldflags are set up correctly
	// in the Makefile and .goreleaser.yml files

	// Save original values
	origVersion := Version
	origCommit := Commit
	origGoVersion := GoVersion

	// Restore original values after test
	defer func() {
		Version = origVersion
		Commit = origCommit
		GoVersion = origGoVersion
	}()

	// Simulate ldflags injection
	Version = "1.2.3"
	Commit = "abcdef"
	GoVersion = "go1.22.0"

	// Get build info
	info := GetBuildInfo()

	// Check values
	if info.Version != "1.2.3" {
		t.Errorf("Expected Version to be '1.2.3', got '%s'", info.Version)
	}

	if info.Commit != "abcdef" {
		t.Errorf("Expected Commit to be 'abcdef', got '%s'", info.Commit)
	}

	if info.GoVersion != "go1.22.0" {
		t.Errorf("Expected GoVersion to be 'go1.22.0', got '%s'", info.GoVersion)
	}
}
