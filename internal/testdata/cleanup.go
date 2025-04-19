package testdata

import (
	"os"
)

// CleanupAllTestData removes all test data files in the given directory.
func CleanupAllTestData(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			os.Remove(dir + "/" + f.Name())
		}
	}
	return nil
}
