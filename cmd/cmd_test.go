package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSubcommandsSetup tests that all subcommands are correctly set up
func TestSubcommandsSetup(t *testing.T) {
	expectedSubcommands := []string{
		"config",
		"credentials",
		"cluster",
		"k3s",
		"kubectl",
		"ssh",
	}
	for _, cmd := range expectedSubcommands {
		subCmd, _, err := RootCmd.Find([]string{cmd})
		assert.NoError(t, err, "Subcommand %s should be found", cmd)
		assert.NotNil(t, subCmd, "Subcommand %s should not be nil", cmd)
	}
	fmt.Println(RootCmd.Version)
}

// TestInvalidSubcommand tests that an error is returned when an invalid subcommand is provided
func TestInvalidSubcommand(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"invalid"}) // Invalid subcommand

	err := RootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `unknown command "invalid" for "h3s"`)
}

// TestPrintWelcome tests that the welcome message is printed when the root command is called without any subcommands
func TestPrintWelcome(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{}) // No arguments to trigger the root command

	err := RootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Welcome to h3s CLI")
	assert.Contains(t, output, "Use --help for more information about available commands")
}

// TestVersionShortFlag tests that the version is printed if the version flag is provided
func TestVersionShortFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"-v"})

	err := RootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version 0.1.0")
}

// TestVersionLongFlag tests that the version is printed if the version flag is provided
func TestVersionLongFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"--version"})

	err := RootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version 0.1.0")
}
