package jsonfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	name = "JSONFile"

	defaultDirPerms = 0750
	defaultFilePerm = 0600
)

type Config struct {
	Dir  string
	File string
}

type JSONFile[T any] struct {
	config   Config
	filename string
}

func New[T any](config Config) (*JSONFile[T], error) {
	err := os.MkdirAll(config.Dir, defaultDirPerms)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	filename, err := createIfNotExist(filepath.Join(config.Dir, config.File))
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
