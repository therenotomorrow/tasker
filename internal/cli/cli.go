package cli

import (
	"context"
	"tasker/internal/usecases"
	"text/template"
)

const (
	success = iota
	failure
	invalid
	unknown
	noArgs

	oneArg  = 1
	twoArgs = 2
)

type Cli struct {
	useCases  *usecases.UseCases
	templates *template.Template
}

func New(storage usecases.Storage) (*Cli, error) {
	templates, err := compileTemplates()
	if err != nil {
		return nil, err
	}

	return &Cli{useCases: usecases.New(storage), templates: templates}, nil
}

func MustNew(storage usecases.Storage) *Cli {
	cli, err := New(storage)
	if err != nil {
		panic(err)
	}

	return cli
}

func (cli *Cli) Dispatch(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.usage()
	}

	command, args := args[0], args[1:]

	return cli.dispatch(ctx, command, args)
}

func (cli *Cli) dispatch(ctx context.Context, command string, args []string) int {
	switch command {
	case "add":
		return cli.add(ctx, args)
	case "update":
		return cli.update(ctx, args)
	case "delete":
		return cli.delete(ctx, args)
	case "mark":
		return cli.mark(ctx, args)
	case "work":
		return cli.work(ctx, args)
	case "done":
		return cli.done(ctx, args)
	case "list":
		return cli.list(ctx, args)
	case "help":
		return cli.help()
	default:
		return cli.unknown(command)
	}
}
