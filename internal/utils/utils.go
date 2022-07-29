package utils

import (
	"fmt"
	"os"
	"path"
)

func MkOutDirIfNotExist(prefix string, dir string) (string, error) {
	newPath := path.Join(prefix, dir)
	info, err := os.Stat(newPath)
	if err != nil {
		err = os.MkdirAll(newPath, 0777)
		if err != nil {
			return "", fmt.Errorf("Failed to create output dir: %s", err)
		}
	} else {
		if !info.IsDir() {
			return "", fmt.Errorf("Output path already exists, but is not a directory.")
		}
	}

	return newPath, nil
}
