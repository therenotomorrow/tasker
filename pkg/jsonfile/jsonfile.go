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

var ErrFileIsNotJSON = errors.New("file is not *.json")

type Config struct {
	File     string
	TestHook func(file *os.File)
}

type JSONFile[T any] struct {
	config   Config
	filename string
}

func New[T any](config Config) (*JSONFile[T], error) {
	file := &JSONFile[T]{config: config, filename: config.File}

	err := file.init()
	if err != nil {
		return nil, err
	}

	return file, nil
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

func (fs *JSONFile[T]) init() error {
	fs.filename = filepath.Clean(fs.filename)

	if filepath.Ext(fs.filename) != ".json" {
		return ErrFileIsNotJSON
	}

	file, err := os.OpenFile(fs.filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, defaultFilePerm)

	switch {
	case errors.Is(err, os.ErrExist):
		return nil
	case err != nil:
		return fmt.Errorf("%s error: %w", name, err)
	}

	defer func() { _ = file.Close() }()

	if fs.config.TestHook != nil {
		fs.config.TestHook(file)
	}

	_, err = file.WriteString("{}")
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	return nil
}
