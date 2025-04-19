package version

import (
	"testing"
)

// TestDefaultVersionValues tests that the default version values are set correctly
func TestDefaultVersionValues(t *testing.T) {
	// Test default version
	if Version != "dev" {
		t.Errorf("Expected default Version to be 'dev', got '%s'", Version)
	}

	// Test default commit
	if Commit != "unknown" {
		t.Errorf("Expected default Commit to be 'unknown', got '%s'", Commit)
	}

	// Test default go version
	if GoVersion != "unknown" {
		t.Errorf("Expected default GoVersion to be 'unknown', got '%s'", GoVersion)
	}
}

// TestGetBuildInfo tests that GetBuildInfo returns the correct values
func TestGetBuildInfo(t *testing.T) {
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

	// Set test values
	Version = "1.0.0"
	Commit = "abc123"
	GoVersion = "go1.22.4"

	// Get build info
	info := GetBuildInfo()

	// Check values
	if info.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got '%s'", info.Version)
	}

	if info.Commit != "abc123" {
		t.Errorf("Expected Commit to be 'abc123', got '%s'", info.Commit)
	}

	if info.GoVersion != "go1.22.4" {
		t.Errorf("Expected GoVersion to be 'go1.22.4', got '%s'", info.GoVersion)
	}
}

// TestBuildInfoStruct tests the BuildInfo struct
func TestBuildInfoStruct(t *testing.T) {
	// Create a BuildInfo struct directly
	info := BuildInfo{
		Version:   "2.0.0",
		Commit:    "def456",
		GoVersion: "go1.23.0",
	}

	// Check values
	if info.Version != "2.0.0" {
		t.Errorf("Expected Version to be '2.0.0', got '%s'", info.Version)
	}

	if info.Commit != "def456" {
		t.Errorf("Expected Commit to be 'def456', got '%s'", info.Commit)
	}

	if info.GoVersion != "go1.23.0" {
		t.Errorf("Expected GoVersion to be 'go1.23.0', got '%s'", info.GoVersion)
	}
}
