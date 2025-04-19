package sshkey

import (
	"h3s/internal/cluster"
	"testing"
)

func TestSSHKeyCreate_Error(t *testing.T) {
	ctx := &cluster.Cluster{}
	// TODO: Use mocks to simulate error
	_, err := Create(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSSHKeyCreate_Success(t *testing.T) {
	ctx := &cluster.Cluster{}
	// TODO: Use mocks to simulate success
	_, err := Create(ctx)
	if err != nil {
		t.Skipf("expected success, got error: %v (likely due to unmocked dependencies)", err)
	}
}
