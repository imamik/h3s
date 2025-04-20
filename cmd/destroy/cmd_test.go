package destroy

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock runDestroyCluster for isolation
type runFunc func(cmd *cobra.Command, args []string) error

func newDestroyCmdWithMock(run runFunc) *cobra.Command {
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy various resources",
		Long:  `Destroy various resources - hetzner cloud resources, k8s components etc.`,
	}
	destroyClusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Destroy an existing cluster",
		Long:  `Destroy an existing cluster including all resources`,
		RunE:  run,
		Args:  cobra.NoArgs,
	}
	destroyCmd.AddCommand(destroyClusterCmd)
	return destroyCmd
}

func TestDestroyCommand_FlagParsing(t *testing.T) {
	mockRun := func(_ *cobra.Command, _ []string) error { return nil }
	buf := new(bytes.Buffer)
	destroyCmd := newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)

	// Valid: destroy cluster --help
	destroyCmd.SetArgs([]string{"cluster", "--help"})
	err := destroyCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Usage:")
	buf.Reset()

	// Invalid: unknown flag
	destroyCmd = newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)
	destroyCmd.SetArgs([]string{"cluster", "--notaflag"})
	err = destroyCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "unknown flag: --notaflag")
	buf.Reset()

	// Invalid: extra argument
	destroyCmd = newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)
	destroyCmd.SetArgs([]string{"cluster", "extraarg"})
	err = destroyCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String()+err.Error(), "unknown command \"extraarg\"")
}

func TestDestroyCommand_ErrorHandling(t *testing.T) {
	errMsg := "mock destroy error"
	mockRun := func(_ *cobra.Command, _ []string) error { return errors.New(errMsg) }
	buf := new(bytes.Buffer)
	destroyCmd := newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)
	destroyCmd.SetArgs([]string{"cluster"})
	err := destroyCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errMsg)
}

func TestDestroyCommand_OutputFormatting(t *testing.T) {
	mockRun := func(_ *cobra.Command, _ []string) error { return nil }
	buf := new(bytes.Buffer)
	destroyCmd := newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)

	// Valid output: cluster help
	destroyCmd.SetArgs([]string{"cluster", "--help"})
	err := destroyCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Usage:")
	buf.Reset()

	// Error output: unknown subcommand
	destroyCmd = newDestroyCmdWithMock(mockRun)
	destroyCmd.SetOut(buf)
	destroyCmd.SetErr(buf)
	destroyCmd.SetArgs([]string{"notacommand"})
	err = destroyCmd.Execute()
	assert.Error(t, err)
	output = buf.String()
	assert.Contains(t, output+err.Error(), "unknown command \"notacommand\"")
}
