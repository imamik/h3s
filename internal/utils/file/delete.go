package file

import (
	"fmt"
	"os"
)

func Delete(filePath string) error {
	if !Exists(filePath) {
		return nil
	}

	absPath, err := Normalize(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Deleting file", absPath)
	err = os.Remove(absPath)
	if err != nil {
		return err
	}

	return nil
}
