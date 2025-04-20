package ssh

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestExecuteWithSSH_Errors(t *testing.T) {
	// nil remote server
	_, err := ExecuteWithSSH("/invalid/path/to/key", nil, "ls")
	if err == nil {
		t.Error("expected error with nil remote server")
	}

	// invalid key path (simulate getConfig failure)
	remote := &hcloud.Server{} // minimal stub
	_, err = ExecuteWithSSH("/invalid/path/to/key", remote, "ls")
	if err == nil {
		t.Error("expected error with invalid key path")
	}
}

func TestExecuteViaProxy_NilGateway(t *testing.T) {
	// Should delegate to ExecuteWithSSH and fail with nil remote
	_, err := ExecuteViaProxy("/invalid/path/to/key", nil, nil, "ls")
	if err == nil {
		t.Error("expected error with nil remote in proxy")
	}
}

func TestExecuteViaProxy_ErrorPaths(t *testing.T) {
	// Simulate error on proxy connection by passing minimal stubs
	gateway := &hcloud.Server{} // minimal stub
	remote := &hcloud.Server{}
	_, err := ExecuteViaProxy("/invalid/path/to/key", gateway, remote, "ls")
	if err == nil {
		t.Error("expected error with invalid proxy connection")
	}
}
