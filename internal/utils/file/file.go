package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// File is an object of a file on the filesystem and holds its most essential data.
// It provides various methods to interact more easy with the file system - read, write, delete and interact with various formats
type File struct {
	path   string
	data   []byte
	errors []error
}

// New creates a new File, normalizes the provided path to its absolute form
func New(path string) *File {
	path = os.ExpandEnv(path)
	absPath, err := filepath.Abs(path)
	f := File{path: absPath}
	if err != nil {
		f.errors = append(f.errors, err)
	}
	return &f
}

// Load reads the file from the filesystem and sets the data of the File
func (f *File) Load() *File {
	data, err := os.ReadFile(f.path)
	if err != nil {
		f.errors = append(f.errors, err)
		return f
	}
	f.SetBytes(data)
	return f
}

// Error returns all errors that have occurred during the lifetime of the File
func (f *File) Error() error {
	// if there are no errors - return nil (no error)
	if len(f.errors) == 0 {
		return nil
	}
	// combine the errors into a single Error
	err := f.errors[0]
	for i := 1; i < len(f.errors); i++ {
		err = fmt.Errorf("%w\n%v", err, f.errors[i])
	}
	return err
}

// Path returns the (absolute) path of the File
func (f *File) Path() (string, error) {
	err := f.Error()
	if err != nil {
		return "", err
	}
	return f.path, nil
}

// Exists checks if the file exists on the filesystem
func (f *File) Exists() bool {
	_, err := os.Stat(f.path)
	return !os.IsNotExist(err)
}

// Save writes the data of the File to the filesystem (creates the file if it does not exist)
func (f *File) Save() (int, error) {
	if f.data == nil {
		return 0, errors.New("no data to save")
	}
	createdFile, err := os.Create(f.path)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			f.errors = append(f.errors, closeErr)
		}
	}(createdFile)

	n, err := createdFile.Write(f.data)
	if err != nil {
		f.errors = append(f.errors, err)
	}

	return n, f.Error()
}

// Delete removes the file from the filesystem and clears the data of the File
func (f *File) Delete() error {
	if !f.Exists() {
		f.SetBytes(nil)
		return nil
	}
	err := os.Remove(f.path)
	if err != nil {
		f.errors = append(f.errors, err)
		return f.Error()
	}
	f.SetBytes(nil)
	return nil
}
