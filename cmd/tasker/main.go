package main

import (
	"cmp"
	"context"
	"os"
	"tasker/internal/cli"
	"tasker/internal/storage"
	"tasker/pkg/jsonfile"
)

func main() {
	ctx := context.Background()

	config := jsonfile.Config{
		Dir:  cmp.Or(os.Getenv("TASKER_DIR"), "."),
		File: "tasker.json",
	}

	tasker := cli.MustNew(storage.MustNew(config))
	status := tasker.Dispatch(ctx, os.Args[1:])

	os.Exit(status)
}
