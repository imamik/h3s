package get

import (
	"bytes"
	"errors"
	"h3s/cmd/dependencies"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// setupTestCmd sets up the command for testing, swapping dependencies for mocks and capturing output.
func setupTestCmd(root *cobra.Command) (*bytes.Buffer, func()) {
	buf := new(bytes.Buffer)
	origOut := root.OutOrStdout()
	root.SetOut(buf)
	// Swap dependencies to mock
	origDeps := dependencies.Get
	dependencies.Get = func() dependencies.CommandDependencies {
		return &dependencies.MockDependencies{
			MockGetClusterContext: func() (*cluster.Cluster, error) { return nil, nil },
			// Add additional mocks as needed
		}
	}
	cleanup := func() {
		root.SetOut(origOut)
		dependencies.Get = origDeps
	}
	return buf, cleanup
}

func TestGetCommand_BaseHelp(t *testing.T) {
	buf, cleanup := setupTestCmd(Cmd)
	defer cleanup()

	Cmd.SetArgs([]string{"--help"})
	err := Cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Load various information about the cluster") {
		t.Errorf("help output missing expected content: %q", out)
	}
}

func TestGetKubeconfig_Success(t *testing.T) {
	_, cleanup := setupTestCmd(Cmd)
	defer cleanup()

	mockCluster := &cluster.Cluster{Config: &config.Config{Domain: "test"}}
	// Mock DownloadKubeconfig to succeed
	dependencies.Get = func() dependencies.CommandDependencies {
		return &dependencies.MockDependencies{
			MockGetClusterContext:  func() (*cluster.Cluster, error) { return mockCluster, nil },
			MockDownloadKubeconfig: func(*cluster.Cluster) error { return nil },
		}
	}

	Cmd.SetArgs([]string{"kubeconfig"})
	err := Cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetKubeconfig_Error(t *testing.T) {
	_, cleanup := setupTestCmd(Cmd)
	defer cleanup()

	mockCluster := &cluster.Cluster{Config: &config.Config{Domain: "test"}}
	dependencies.Get = func() dependencies.CommandDependencies {
		return &dependencies.MockDependencies{
			MockGetClusterContext:  func() (*cluster.Cluster, error) { return mockCluster, nil },
			MockDownloadKubeconfig: func(*cluster.Cluster) error { return errors.New("fail download") },
		}
	}

	Cmd.SetArgs([]string{"kubeconfig"})
	err := Cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "fail download") {
		t.Errorf("expected error for kubeconfig download failure, got: %v", err)
	}
}

func TestGetToken_Success(t *testing.T) {
	_, cleanup := setupTestCmd(Cmd)
	defer cleanup()

	mockCluster := &cluster.Cluster{Config: &config.Config{Domain: "test"}}
	dependencies.Get = func() dependencies.CommandDependencies {
		return &dependencies.MockDependencies{
			MockGetClusterContext:   func() (*cluster.Cluster, error) { return mockCluster, nil },
			MockGenerateK8sToken:    func(*cluster.Cluster, string, string, int) (string, error) { return "token123", nil },
			MockExecuteLocalCommand: func(string) (string, error) { return "", nil },
		}
	}
	Cmd.SetArgs([]string{"token"})
	err := Cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetToken_Error(t *testing.T) {
	_, cleanup := setupTestCmd(Cmd)
	defer cleanup()

	mockCluster := &cluster.Cluster{Config: &config.Config{Domain: "test"}}
	dependencies.Get = func() dependencies.CommandDependencies {
		return &dependencies.MockDependencies{
			MockGetClusterContext: func() (*cluster.Cluster, error) { return mockCluster, nil },
			MockGenerateK8sToken:  func(*cluster.Cluster, string, string, int) (string, error) { return "", errors.New("fail token") },
		}
	}
	Cmd.SetArgs([]string{"token"})
	err := Cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "fail token") {
		t.Errorf("expected error for token failure, got: %v", err)
	}
}

func TestGet_InvalidCommand(t *testing.T) {
	_, cleanup := setupTestCmd(Cmd)
	defer cleanup()
	Cmd.SetArgs([]string{"doesnotexist"})
	err := Cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("expected error for unknown command, got: %v", err)
	}
}
