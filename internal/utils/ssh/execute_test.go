package ssh

import (
	"errors"
	"testing"
)

func TestRemoveKnownHostsEntry_InvalidPath(t *testing.T) {
	err := removeKnownHostsEntry("") // Empty IP, should not match any line, but test for unexpected errors
	if err != nil && !errors.Is(err, nil) {
		t.Errorf("Expected no error or nil error, got: %v", err)
	}
}
