package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAndPath(t *testing.T) {
	f := New("testfile.txt")
	path, err := f.Path()
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
}

func TestExistsAndSaveAndDelete(t *testing.T) {
	filename := "testfile_exists.txt"
	f := New(filename)
	f.SetString("data")
	_, err := f.Save()
	assert.NoError(t, err)
	assert.True(t, f.Exists())
	err = f.Delete()
	assert.NoError(t, err)
	assert.False(t, f.Exists())
}

func TestErrorAggregation(t *testing.T) {
	f := &File{}
	f.errors = append(f.errors, os.ErrNotExist, os.ErrPermission)
	err := f.Error()
	assert.Error(t, err)
	assert.NotEmpty(t, err.Error())
}

func TestLoadNonExistentFile(t *testing.T) {
	f := New("nonexistent_file.txt")
	f.Load()
	assert.Error(t, f.Error())
}

func TestSaveWithNilData(t *testing.T) {
	f := New("should_not_create.txt")
	_, err := f.Save()
	assert.Error(t, err)
}

func TestDeleteNonExistentFile(t *testing.T) {
	f := New("definitely_not_exists.txt")
	err := f.Delete()
	assert.NoError(t, err)
}

func TestSaveAndLoadEmptyFile(t *testing.T) {
	filename := "emptyfile.txt"
	f := New(filename)
	f.SetBytes([]byte{})
	_, err := f.Save()
	assert.NoError(t, err)
	f2 := New(filename)
	f2.Load()
	assert.Empty(t, string(f2.data))
	_ = f.Delete()
}

func TestPermissionDenied(t *testing.T) {
	filename := "readonly.txt"
	f := New(filename)
	f.SetString("data")
	_, err := f.Save()
	assert.NoError(t, err)

	// Create file and set to read-only
	_, err = f.Save()
	assert.NoError(t, err)

	// Change permissions to read-only
	err = os.Chmod(filename, 0o400)
	assert.NoError(t, err)

	// Attempt to overwrite read-only file
	f.SetString("newdata")
	_, err = f.Save()
	assert.Error(t, err)

	// Restore permissions to allow deletion
	err = os.Chmod(filename, 0o600)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	f2 := New(filename)
	f2.Load()
	assert.Len(t, f2.data, len(largeData))
	_ = f.Delete()
}
