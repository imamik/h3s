package file

import (
	"os"
	"testing"
)

func TestNewAndPath(t *testing.T) {
	f := New("testfile.txt")
	path, err := f.Path()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Errorf("expected non-empty path")
	}
}

func TestExistsAndSaveAndDelete(t *testing.T) {
	filename := "testfile_exists.txt"
	f := New(filename)
	f.SetString("data")
	_, err := f.Save()
	if err != nil {
		t.Fatalf("unexpected error on save: %v", err)
	}
	if !f.Exists() {
		t.Errorf("file should exist after save")
	}
	if err := f.Delete(); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}
	if f.Exists() {
		t.Errorf("file should not exist after delete")
	}
}

func TestErrorAggregation(t *testing.T) {
	f := &File{}
	f.errors = append(f.errors, os.ErrNotExist, os.ErrPermission)
	err := f.Error()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	f := New("nonexistent_file.txt")
	f.Load()
	if err := f.Error(); err == nil {
		t.Errorf("expected error when loading non-existent file, got nil")
	}
}

func TestSaveWithNilData(t *testing.T) {
	f := New("should_not_create.txt")
	_, err := f.Save()
	if err == nil {
		t.Errorf("expected error when saving with nil data, got nil")
	}
}

func TestDeleteNonExistentFile(t *testing.T) {
	f := New("definitely_not_exists.txt")
	err := f.Delete()
	if err != nil {
		t.Errorf("expected no error when deleting non-existent file, got %v", err)
	}
}

func TestSaveAndLoadEmptyFile(t *testing.T) {
	filename := "emptyfile.txt"
	f := New(filename)
	f.SetBytes([]byte{})
	_, err := f.Save()
	if err != nil {
		t.Fatalf("unexpected error on save: %v", err)
	}
	f2 := New(filename)
	f2.Load()
	if string(f2.data) != "" {
		t.Errorf("expected empty file content, got %q", string(f2.data))
	}
	_ = f.Delete()
}

func TestPermissionDenied(t *testing.T) {
	filename := "readonly.txt"
	f := New(filename)
	f.SetString("data")
	// Create file and set to read-only
	_, err := f.Save()
	if err != nil {
		t.Fatalf("unexpected error on save: %v", err)
	}
	err = os.Chmod(filename, 0400)
	if err != nil {
		t.Fatalf("failed to set file to read-only: %v", err)
	}
	// Attempt to overwrite read-only file
	f.SetString("newdata")
	_, err = f.Save()
	if err == nil {
		t.Errorf("expected permission error when saving to read-only file, got nil")
	}
	_ = os.Chmod(filename, 0600)
	_ = f.Delete()
}

func TestLargeFile(t *testing.T) {
	filename := "largefile.txt"
	f := New(filename)
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = 'A'
	}
	f.SetBytes(largeData)
	_, err := f.Save()
	if err != nil {
		t.Fatalf("unexpected error on save: %v", err)
	}
	f2 := New(filename)
	f2.Load()
	if len(f2.data) != len(largeData) {
		t.Errorf("expected %d bytes, got %d", len(largeData), len(f2.data))
	}
	_ = f.Delete()
}
