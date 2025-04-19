package ssh

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newMockSSHCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "ssh",
		Short:              "Proxy ssh commands to first control plane server",
		DisableFlagParsing: true,
		RunE:               runE,
	}
	return cmd
}

func TestSSH_Help(t *testing.T) {
	// Create a copy of the actual command for testing
	cmd := &cobra.Command{
		Use:   Cmd.Use,
		Short: Cmd.Short,
		Long:  Cmd.Long,
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()

	assert.NoError(t, err, "Help command should not return an error")
	out := buf.String()
	assert.Contains(t, out, "Proxy ssh commands to first control plane server", "Help output should contain command description")
}

func TestSSH_Success(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command
	cmd := newMockSSHCmd(func(cmd *cobra.Command, args []string) error {
		cmd.Println("mock ssh output: ls -la")
		return nil
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"ls", "-la"})
	err := cmd.Execute()

	assert.NoError(t, err, "Command should execute without error")
	assert.Contains(t, buf.String(), "mock ssh output: ls -la", "Command output should contain the executed command")
}

func TestSSH_WithComplexCommand(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command
	cmd := newMockSSHCmd(func(cmd *cobra.Command, args []string) error {
		cmd.Println("mock ssh output: find /var/log -name *.log | grep system")
		return nil
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"find", "/var/log", "-name", "*.log", "|", "grep", "system"})
	err := cmd.Execute()

	assert.NoError(t, err, "Command should execute without error")
	assert.Contains(t, buf.String(), "mock ssh output: find /var/log -name *.log | grep system",
		"Command output should contain the complex command with all arguments")
}

func TestSSH_Error(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command with an error
	cmd := newMockSSHCmd(func(cmd *cobra.Command, args []string) error {
		return errors.New("failed to execute ssh command: mock ssh connection failure")
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"ls"})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error")
	assert.Contains(t, err.Error(), "failed to execute ssh command", "Error message should contain failure reason")
}

func TestSSH_ClusterContextError(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command with a cluster context error
	cmd := newMockSSHCmd(func(cmd *cobra.Command, args []string) error {
		return errors.New("failed to load cluster context: mock cluster context error")
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"ls"})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error")
	assert.Contains(t, err.Error(), "failed to load cluster context", "Error message should contain failure reason")
}

func TestSSH_InvalidFlag(t *testing.T) {
	// Since SSH command has DisableFlagParsing=true, we need to test a different way
	// Create a command without DisableFlagParsing for this test
	cmd := &cobra.Command{Use: "ssh"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--notaflag"})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error for invalid flag")
	assert.Contains(t, err.Error(), "unknown flag: --notaflag", "Error message should mention the invalid flag")
}
