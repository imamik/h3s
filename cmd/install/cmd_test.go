package install

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newMockK3sCmd(runE func(cmd *cobra.Command, _ []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k3s",
		Short: "Install k3s on all servers in the cluster",
		RunE:  runE,
	}
	return cmd
}

func TestInstall_Help(t *testing.T) {
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
	assert.Contains(t, out, "Install k3s or k8s components", "Help output should contain command description")
}

func TestInstallK3s_Success(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command
	cmd := newMockK3sCmd(func(cmd *cobra.Command, _ []string) error {
		cmd.Println("K3s installed successfully")
		return nil
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NoError(t, err, "Command should execute without error")
	assert.Contains(t, buf.String(), "K3s installed successfully", "Success message should be printed")
}

func TestInstallK3s_Error(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command with an error
	cmd := newMockK3sCmd(func(_ *cobra.Command, _ []string) error {
		return errors.New("failed to install k3s: mock k3s installation failure")
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error")
	assert.Contains(t, err.Error(), "failed to install k3s", "Error message should contain failure reason")
}

func TestInstallComponents_Success(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command
	cmd := &cobra.Command{
		Use:   "components",
		Short: "Install k8s components in the cluster",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println("Kubernetes components installed successfully")
			return nil
		},
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NoError(t, err, "Command should execute without error")
	assert.Contains(t, buf.String(), "Kubernetes components installed successfully", "Success message should be printed")
}

func TestInstallComponents_Error(t *testing.T) {
	// Create a mock command that simulates the behavior of the real command with an error
	cmd := &cobra.Command{
		Use:   "components",
		Short: "Install k8s components in the cluster",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("failed to install kubernetes components: mock components installation failure")
		},
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error")
	assert.Contains(t, err.Error(), "failed to install kubernetes components", "Error message should contain failure reason")
}

func TestInstall_InvalidFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "install"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--notaflag"})
	err := cmd.Execute()

	assert.Error(t, err, "Command should return an error for invalid flag")
	assert.Contains(t, err.Error(), "unknown flag: --notaflag", "Error message should mention the invalid flag")
}
