package version

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLdflagsInBuildConfigs tests that the ldflags are correctly set up in the Makefile and .goreleaser.yml files
func TestLdflagsInBuildConfigs(t *testing.T) {
	// Find the project root directory
	projectRoot := findProjectRoot(t)
	if projectRoot == "" {
		t.Skip("Could not find project root directory")
	}

	// Check Makefile
	t.Run("Makefile", func(t *testing.T) {
		makefilePath := filepath.Join(projectRoot, "Makefile")
		content, err := os.ReadFile(makefilePath)
		if err != nil {
			t.Fatalf("Failed to read Makefile: %v", err)
		}

		makefileContent := string(content)

		// Check that the ldflags are set up correctly
		requiredLdflags := []string{
			"-X h3s/internal/version.Version=$(VERSION)",
			"-X h3s/internal/version.Commit=$(shell git rev-parse --short HEAD",
			"-X 'h3s/internal/version.GoVersion=$(shell go version",
		}

		for _, flag := range requiredLdflags {
			if !strings.Contains(makefileContent, flag) {
				t.Errorf("Makefile does not contain required ldflags: %s", flag)
			}
		}
	})

	// Check .goreleaser.yml
	t.Run(".goreleaser.yml", func(t *testing.T) {
		goreleaserPath := filepath.Join(projectRoot, ".goreleaser.yml")
		content, err := os.ReadFile(goreleaserPath)
		if err != nil {
			t.Fatalf("Failed to read .goreleaser.yml: %v", err)
		}

		goreleaserContent := string(content)

		// Check that the ldflags are set up correctly
		requiredLdflags := []string{
			"-X h3s/internal/version.Version={{ .Version }}",
			"-X h3s/internal/version.Commit={{ .Commit }}",
			"-X 'h3s/internal/version.GoVersion={{ .Env.GOVERSION }}'",
		}

		for _, flag := range requiredLdflags {
			if !strings.Contains(goreleaserContent, flag) {
				t.Errorf(".goreleaser.yml does not contain required ldflags: %s", flag)
			}
		}
	})
}

// findProjectRoot finds the root directory of the project
func findProjectRoot(t *testing.T) string {
	// Start from the current directory and go up until we find a go.mod file
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory
			return ""
		}
		dir = parent
	}
}
