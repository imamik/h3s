package ping

import (
	"testing"
	"time"
)

func TestIsServerAvailable_InvalidIP(t *testing.T) {
	if isServerAvailable("") {
		t.Error("expected false for empty IP")
	}
	if isServerAvailable("256.256.256.256") {
		t.Error("expected false for invalid IP")
	}
}

func TestPing_NilServer(t *testing.T) {
	// Should not panic, but will likely cause a nil pointer dereference if not handled
	defer func() {
		if r := recover(); r == nil {
			t.Error("Ping(nil, ...) should panic on nil server (for now)")
		}
	}()
	Ping(nil, 1*time.Millisecond)
}
