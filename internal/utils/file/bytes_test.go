package file

import "testing"

func TestSetAndGetBytes(t *testing.T) {
	f := &File{}
	input := []byte{1, 2, 3, 4, 5}
	f.SetBytes(input)
	out, err := f.GetBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(input) {
		t.Errorf("GetBytes() length = %d, want %d", len(out), len(input))
	}
	for i := range input {
		if out[i] != input[i] {
			t.Errorf("GetBytes()[%d] = %d, want %d", i, out[i], input[i])
		}
	}
}

func TestGetBytesWithError(t *testing.T) {
	f := &File{}
	f.errors = append(f.errors, errFake{})
	_, err := f.GetBytes()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
