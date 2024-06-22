package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Normalize(filePath string) (string, error) {
	filePath = os.ExpandEnv(filePath)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return filePath, fmt.Errorf("error resolving absolute path of %s: %w", filePath, err)
	}
	return absPath, nil
}
