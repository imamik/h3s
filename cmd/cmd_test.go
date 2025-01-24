package cmd

import (
	"bytes"
	"fmt"
	"h3s/internal/version"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSubcommandsSetup tests that all subcommands are correctly set up
func TestSubcommandsSetup(t *testing.T) {
	Initialize(version.BuildInfo{})
	expectedSubcommands := []string{
		"create",
		"destroy",
		"get",
		"install",
		"kubectl",
		"ssh",
	}
	for _, cmd := range expectedSubcommands {
		subCmd, _, err := Cmd.Find([]string{cmd})
		assert.NoError(t, err, "Subcommand %s should be found", cmd)
		assert.NotNil(t, subCmd, "Subcommand %s should not be nil", cmd)
	}
	fmt.Println(Cmd.Version)
}

// TestInvalidSubcommand tests that an error is returned when an invalid subcommand is provided
func TestInvalidSubcommand(t *testing.T) {
	Initialize(version.BuildInfo{})
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs([]string{"invalid"}) // Invalid subcommand

	err := Cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `unknown command "invalid" for "h3s"`)
}

// TestPrintWelcome tests that the welcome message is printed when the root command is called without any subcommands
func TestPrintWelcome(t *testing.T) {
	Initialize(version.BuildInfo{})
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{})

	err := Cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Welcome to h3s CLI")
	assert.Contains(t, output, "Use --help for more information about available commands")
}

// TestVersionShortFlag tests that the version is printed if the version flag is provided
func TestVersionShortFlag(t *testing.T) {
	Initialize(version.BuildInfo{Version: "dev"})
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{"-v"})

	err := Cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version dev")
}

// TestVersionLongFlag tests that the version is printed if the version flag is provided
func TestVersionLongFlag(t *testing.T) {
	Initialize(version.BuildInfo{Version: "dev"})
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{"--version"})

	err := Cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version dev")
}

// TestVersionOutput tests that the version output includes all necessary information
func TestVersionOutput(t *testing.T) {
	Initialize(version.BuildInfo{
		Version:   "dev",
		Commit:    "test-commit",
		GoVersion: "go1.21",
	})

	var buf bytes.Buffer
	Cmd.SetOut(&buf)
	Cmd.SetArgs([]string{"--version"})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	expected := []string{
		"h3s version dev",
		"Commit: test-commit",
		"Go version: go1.21",
	}

	for _, s := range expected {
		if !strings.Contains(output, s) {
			t.Errorf("Expected %q in version output", s)
		}
	}
}
