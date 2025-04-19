package version

import (
	"fmt"
	"strings"
	"testing"
)

// TestVersionStringFormatting tests the formatting of version strings
func TestVersionStringFormatting(t *testing.T) {
	// Test cases for version formatting
	testCases := []struct {
		version   string
		commit    string
		goVersion string
		expected  string
	}{
		{
			version:   "1.0.0",
			commit:    "abc123",
			goVersion: "go1.22.4",
			expected:  "1.0.0\nCommit: abc123\nGo version: go1.22.4",
		},
		{
			version:   "v2.1.0",
			commit:    "def456",
			goVersion: "go1.23.0",
			expected:  "v2.1.0\nCommit: def456\nGo version: go1.23.0",
		},
		{
			version:   "dev",
			commit:    "unknown",
			goVersion: "unknown",
			expected:  "dev\nCommit: unknown\nGo version: unknown",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			// Create a BuildInfo struct with the test case values
			info := BuildInfo{
				Version:   tc.version,
				Commit:    tc.commit,
				GoVersion: tc.goVersion,
			}

			// Format the version string
			formatted := fmt.Sprintf("%s\nCommit: %s\nGo version: %s",
				info.Version,
				info.Commit,
				info.GoVersion,
			)

			// Check that the formatted string matches the expected output
			if formatted != tc.expected {
				t.Errorf("Expected formatted string to be '%s', got '%s'", tc.expected, formatted)
			}
		})
	}
}

// TestVersionStringParsing tests parsing of version strings
func TestVersionStringParsing(t *testing.T) {
	// Test cases for version parsing
	testCases := []struct {
		input         string
		expectVersion string
		expectCommit  string
		expectGo      string
	}{
		{
			input:         "1.0.0\nCommit: abc123\nGo version: go1.22.4",
			expectVersion: "1.0.0",
			expectCommit:  "abc123",
			expectGo:      "go1.22.4",
		},
		{
			input:         "v2.1.0\nCommit: def456\nGo version: go1.23.0",
			expectVersion: "v2.1.0",
			expectCommit:  "def456",
			expectGo:      "go1.23.0",
		},
		{
			input:         "dev\nCommit: unknown\nGo version: unknown",
			expectVersion: "dev",
			expectCommit:  "unknown",
			expectGo:      "unknown",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			// Parse the input string
			lines := strings.Split(tc.input, "\n")
			if len(lines) < 3 {
				t.Fatalf("Input string does not have enough lines: %s", tc.input)
			}

			// Extract version, commit, and go version
			version := lines[0]
			commit := strings.TrimPrefix(lines[1], "Commit: ")
			goVersion := strings.TrimPrefix(lines[2], "Go version: ")

			// Check that the parsed values match the expected values
			if version != tc.expectVersion {
				t.Errorf("Expected version to be '%s', got '%s'", tc.expectVersion, version)
			}
			if commit != tc.expectCommit {
				t.Errorf("Expected commit to be '%s', got '%s'", tc.expectCommit, commit)
			}
			if goVersion != tc.expectGo {
				t.Errorf("Expected go version to be '%s', got '%s'", tc.expectGo, goVersion)
			}
		})
	}
}
