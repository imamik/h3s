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
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{"-v"})

	err := Cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version 0.1.0")
}

// TestVersionLongFlag tests that the version is printed if the version flag is provided
func TestVersionLongFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{"--version"})

	err := Cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "h3s version 0.1.0")
}
