package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/pkg/utils/file"
	"os"
)

func Save(config any, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config data: %w", err)
	}

	// Convert to absolute path
	absPath, err := file.Normalize(filePath)
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
