package jsonfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	name = "JSONFile"

	defaultFilePerm = 0o600
)

type Config struct {
	File string
}

type JSONFile[T any] struct {
	config   Config
	filename string
}

func New[T any](config Config) (*JSONFile[T], error) {
	filename, err := createIfNotExist(config.File)
	if err != nil {
		return nil, err
	}

	return &JSONFile[T]{config: config, filename: filename}, nil
}

func (fs *JSONFile[T]) Load() (T, error) {
	var data T

	file, err := os.Open(fs.filename)
	if err != nil {
		return data, fmt.Errorf("%s load error: %w", name, err)
	}

	defer func() { _ = file.Close() }()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("%s decode error: %w", name, err)
	}

	return data, nil
}

func (fs *JSONFile[T]) Save(data T) error {
	file, err := os.Create(fs.filename)
	if err != nil {
		return fmt.Errorf("%s save error: %w", name, err)
	}

	defer func() { _ = file.Close() }()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("%s encode error: %w", name, err)
	}

	return nil
}

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
