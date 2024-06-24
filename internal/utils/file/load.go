package file

import (
	"fmt"
	"os"
)

func Load(filePath string) ([]byte, error) {
	absPath, err := Normalize(filePath)
	if err != nil {
		return nil, fmt.Errorf("error normalizing file path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", absPath, err)
	}

	return data, err
}
