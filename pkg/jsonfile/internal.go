package jsonfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func createIfNotExist(filename string) (string, error) {
	filename = filepath.Clean(filename)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, defaultFilePerm)

	switch {
	case errors.Is(err, os.ErrExist):
		return filename, nil
	case err != nil:
		return "", fmt.Errorf("%s error: %w", name, err)
	}

	defer func() { _ = file.Close() }()

	_, err = file.WriteString("{}")
	if err != nil {
		return "", fmt.Errorf("%s error: %w", name, err)
	}

	return filename, nil
}
