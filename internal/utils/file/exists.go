package file

import (
	"fmt"
	"os"
)

func Exists(filePath string) bool {
	absPath, err := Normalize(filePath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = os.Stat(absPath)
	return !os.IsNotExist(err)
}
