package file

import (
	"testing"
)

func TestSetAndGetString(t *testing.T) {
	f := &File{}
	input := "hello world"
	f.SetString(input)
	out, err := f.GetString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("GetString() = %q, want %q", out, input)
	}
}

func TestGetStringWithError(t *testing.T) {
	f := &File{}
	f.errors = append(f.errors, errFake{})
	_, err := f.GetString()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestSetAndGetEmptyString(t *testing.T) {
	f := &File{}
	f.SetString("")
	out, err := f.GetString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("GetString() = %q, want empty string", out)
	}
}

func TestSetAndGetLargeString(t *testing.T) {
	f := &File{}
	large := make([]byte, 100000)
	for i := range large {
		large[i] = 'x'
	}
	f.SetString(string(large))
	out, err := f.GetString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 100000 {
		t.Errorf("GetString() length = %d, want 100000", len(out))
	}
}

type errFake struct{}

func (e errFake) Error() string { return "fake error" }
