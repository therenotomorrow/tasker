package main

import (
	"cmp"
	"context"
	"os"

	"github.com/therenotomorrow/tasker/internal/cli"
	"github.com/therenotomorrow/tasker/internal/storage"
	"github.com/therenotomorrow/tasker/pkg/jsonfile"
)

func main() {
	ctx := context.Background()
	file := cmp.Or(os.Getenv("TASKER_FILE"), "tasker.json")

	config := jsonfile.Config{File: file}
	tasker := cli.MustNew(storage.MustNew(config))
	status := tasker.Dispatch(ctx, os.Args[1:])

	os.Exit(status)
}
