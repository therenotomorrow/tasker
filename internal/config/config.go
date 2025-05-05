package config

import (
	"path/filepath"
	"runtime"
)

func root() string {
	_, filename, _, ok := runtime.Caller(0)
	_ = ok

	filename = filepath.Join(filepath.Dir(filename), "..", "..")

	return filepath.Clean(filename)
}

func Path(parts ...string) string {
	parts = append([]string{root()}, parts...)

	return filepath.Clean(filepath.Join(parts...))
}
