package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// runMainWithArgs runs the main function with the provided args and checks if the output contains the expected strings.
func runMainWithArgs(args []string, expectedOutput []string, t *testing.T) {
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
	go func() {
		cobra.MousetrapHelpText = "" // Disable mousetraps to prevent extra output on Windows
		main()
		wOut.Close()
		wErr.Close()
	}()

	// Read the captured output
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

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

func TestMainOutput(t *testing.T) {
	runMainWithArgs(
		[]string{"h3s"},
		[]string{
			"Welcome to h3s CLI",
			"Use --help for more information about available commands",
		},
		t)
}

func TestMainHelpOutput(t *testing.T) {
	runMainWithArgs(
		[]string{"h3s", "--help"},
		[]string{
			"config",
			"credentials",
			"cluster",
			"k3s",
			"kubectl",
			"ssh",
		},
		t)
}

func TestMainInvalidArgOutput(t *testing.T) {
	runMainWithArgs(
		[]string{"h3s", "invalid"},
		[]string{
			"Error: unknown command \"invalid\" for \"h3s\"",
		},
		t)
}