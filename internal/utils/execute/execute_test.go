package execute

import (
	"strings"
	"testing"
)

func TestLocal_Success(t *testing.T) {
	out, err := Local("echo hello")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected output to contain 'hello', got: %q", out)
	}
}

func TestLocal_Failure(t *testing.T) {
	_, err := Local("nonexistentcommand1234")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}
