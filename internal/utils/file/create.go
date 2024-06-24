package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Create(filePath string) error {
	if Exists(filePath) {
		return nil
	}

	absPath, err := Normalize(filePath)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(absPath)
	if !Exists(dirPath) {
		err := os.MkdirAll(dirPath, 0755) // 0755 is commonly used permission for directories
		if err != nil {
			return err
		}
	}

	fmt.Println("Creating file", absPath)
	file, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	return nil
}
