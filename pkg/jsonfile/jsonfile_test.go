package jsonfile_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/therenotomorrow/tasker/internal/config"
	"github.com/therenotomorrow/tasker/pkg/jsonfile"
)

const (
	strictPerms  = 0o100
	regularPerms = 0o600
)

type Type struct {
	Str string `json:"str"`
	Num int    `json:"num"`
}

func TestIntegrationNewNotJSON(t *testing.T) {
	t.Parallel()

	got, err := jsonfile.New[Type](jsonfile.Config{File: t.Name()})

	if !errors.Is(err, jsonfile.ErrFileIsNotJSON) {
		t.Fatalf("New() got = %v, error = %v, want = %v", got, err, jsonfile.ErrFileIsNotJSON)
	}
}

func TestIntegrationNewExist(t *testing.T) {
	t.Parallel()

	filename := t.Name() + ".jsosn"

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filepath.Clean(filename))

	fmt.Println("sad")

	got, err := jsonfile.New[Type](jsonfile.Config{File: filename})

	if err != nil || got == nil {
		t.Fatalf("New() got = %v, error = %v, want = %v", got, err, nil)
	}
}

func TestIntegrationNewInvalidFilename(t *testing.T) {
	t.Parallel()

	got, err := jsonfile.New[Type](jsonfile.Config{File: string([]byte{0}) + ".json"})

	if want := "invalid argument"; got != nil || err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("New() got = %v, error = %v, want = %v", got, err, want)
	}
}

func TestIntegrationNewCannotWrite(t *testing.T) {
	t.Parallel()

	filename := t.Name() + ".json"

	defer func() { _ = os.Remove(filename) }()

	cfg := jsonfile.Config{File: filename, TestHook: func(file *os.File) {
		_ = file.Close()
	}}
	got, err := jsonfile.New[Type](cfg)

	if want := "file already closed"; got != nil || err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("New() got = %v, error = %v, want = %v", got, err, want)
	}
}

func TestIntegrationJSONFileLoad(t *testing.T) {
	t.Parallel()

	filename := t.Name() + ".json"

	defer func() { _ = os.Remove(filename) }()

	_ = os.WriteFile(filename, []byte("{}"), strictPerms)

	file, _ := jsonfile.New[Type](jsonfile.Config{File: filename})

	got, err := file.Load()
	if err == nil {
		t.Fatalf("Load() got = %v, error = %v, want = %v", got, err, nil)
	}

	_ = os.Remove(filename)

	file, _ = jsonfile.New[Type](jsonfile.Config{File: config.Path("test", "data", "empty.json")})

	got, err = file.Load()
	if err == nil {
		t.Fatalf("Load() got = %v, error = %v, want = %v", got, err, nil)
	}

	_ = os.WriteFile(filename, []byte(`{"str":"24","num":42}`), regularPerms)

	file, _ = jsonfile.New[Type](jsonfile.Config{File: filename})
	got, err = file.Load()
	want := Type{Str: "24", Num: 42}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Load() got = %v, error = %v, want = %v", got, err, want)
	}
}

func TestIntegrationJSONFileSave(t *testing.T) {
	t.Parallel()

	data := Type{Str: "42", Num: 24}
	filename := t.Name() + ".json"

	defer func() { _ = os.Remove(filename) }()

	_ = os.WriteFile(filename, []byte("{}"), strictPerms)

	file, _ := jsonfile.New[any](jsonfile.Config{File: filename})

	err := file.Save(data)
	if err == nil {
		t.Fatalf("Save() error = %v, want = %v", err, nil)
	}

	_ = os.Remove(filename)

	file, _ = jsonfile.New[any](jsonfile.Config{File: filename})

	err = file.Save(func() {})
	if err == nil {
		t.Fatalf("Save() error = %v, want = %v", err, nil)
	}

	err = file.Save(data)
	if err != nil {
		t.Fatalf("Save() error = %v, want = %v", err, nil)
	}

	got, _ := os.ReadFile(filepath.Clean(filename))
	want := `{
  "str": "42",
  "num": 24
}
`

	if got := string(got); !reflect.DeepEqual(got, want) {
		t.Fatalf("Save() got = %v, error = %v, want = %v", got, err, want)
	}
}
