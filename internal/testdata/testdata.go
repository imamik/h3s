// Package testdata provides utilities for managing test data files.
package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const Version = "v1"

// LoadTestData loads test data from a file and unmarshals it into v.
func LoadTestData(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// SaveTestData saves the given object as test data to a file.
func SaveTestData(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// CleanupTestData removes the specified test data file.
func CleanupTestData(path string) error {
	return os.Remove(path)
}

// GenerateTestConfig returns a config struct with the given properties.
type TestConfig struct {
	Name    string `json:"name"`
	Valid   bool   `json:"valid"`
	Version string `json:"version"`
}

func GenerateTestConfig(name string, valid bool) TestConfig {
	return TestConfig{
		Name:    name,
		Valid:   valid,
		Version: Version,
	}
}

// ListTestDataFiles returns all test data files in the directory.
func ListTestDataFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
