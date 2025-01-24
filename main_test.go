package main

import (
	"bytes"
	"h3s/cmd"
	"h3s/internal/version"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// runMainWithArgs runs the main function with the provided args and checks if the output contains the expected strings.
func runMainWithArgs(args, expectedOutput []string, t *testing.T) {
	// Save the original stdout, stderr, and os.Args
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	originalArgs := os.Args

	// Create pipes to capture stdout and stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Override os.Args to simulate running the program with the provided arguments
	os.Args = args

	// Run the main function
	done := make(chan struct{})
	go func() {
		defer close(done)
		cobra.MousetrapHelpText = "" // Disable mousetraps to prevent extra output on Windows
		main()
		wOut.Close()
		wErr.Close()
	}()

	// Wait for the main function to complete
	<-done

	// Read the captured output
	var bufOut, bufErr bytes.Buffer
	_, err := io.Copy(&bufOut, rOut)
	if err != nil {
		t.Fatalf("Failed to copy stdout: %v", err)
	}
	_, err = io.Copy(&bufErr, rErr)
	if err != nil {
		t.Fatalf("Failed to copy stderr: %v", err)
	}

	// Restore original stdout, stderr, and os.Args
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	os.Args = originalArgs

	// Combine stdout and stderr
	output := bufOut.String() + bufErr.String()

	// Check for expected output
	for _, str := range expectedOutput {
		if !strings.Contains(output, str) {
			t.Errorf("Expected output to contain %q but it did not", str)
		}
	}
}

// TestMainOutput tests the output when no arguments are provided - it should display the welcome message and "--help" hint
func TestMainOutput(t *testing.T) {
	runMainWithArgs(
		[]string{"h3s"},
		[]string{
			"Welcome to h3s CLI",
			"Use --help for more information about available commands",
		},
		t)
}

// TestMainHelpOutput tests the output when the help flag is provided - it should list all available commands
func TestMainHelpOutput(t *testing.T) {
	runMainWithArgs(
		[]string{"h3s", "--help"},
		[]string{
			"cluster",
			"k3s",
			"kubectl",
			"ssh",
		},
		t)
}

// TestMainInvalidArgOutput tests the output when an invalid argument is provided.
func TestMainInvalidArgOutput(t *testing.T) {
	// Create buffers to capture output
	var stdout, stderr bytes.Buffer

	// Create a new cobra command for testing
	cmd := cmd.Cmd
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"invalid"})

	// Execute the command
	err := cmd.Execute()

	// Print the actual output for debugging
	t.Logf("Actual stdout: %q", stdout.String())
	t.Logf("Actual stderr: %q", stderr.String())

	// Verify error occurred
	if err == nil {
		t.Error("Expected an error but got none")
	}

	// Expected error message
	expectedError := "Error: unknown command \"invalid\" for \"h3s\""
	if !strings.Contains(stderr.String(), expectedError) {
		t.Errorf("Expected stderr to contain %q but it did not", expectedError)
	}

	// Expected help message
	expectedHelp := "Run 'h3s --help' for usage."
	if !strings.Contains(stderr.String(), expectedHelp) {
		t.Errorf("Expected stderr to contain %q but it did not", expectedHelp)
	}
}

// TestVersionInjection tests the version injection via ldflags
func TestVersionInjection(t *testing.T) {
	// Test default version
	if version.Version != "dev" {
		t.Errorf("Expected default version 'dev', got '%s'", version.Version)
	}

	// Test build-time injection would require building with ldflags,
	// which is better handled in integration tests
}
