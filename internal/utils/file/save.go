package file

import (
	"fmt"
	"os"
)

func Save(data []byte, filePath string) error {
	// Convert to absolute path
	absPath, err := Normalize(filePath)
	if err != nil {
		return fmt.Errorf("error resolving absolute path of %s: %w", filePath, err)
	}

	createdFile, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", absPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing createdFile %s: %s", absPath, err)
		}
	}(createdFile)

	_, err = createdFile.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", absPath, err)
	}

	return nil
}
