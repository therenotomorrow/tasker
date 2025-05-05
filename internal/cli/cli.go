package cli

import (
	"context"
	"text/template"

	"github.com/therenotomorrow/tasker/internal/usecases"
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
	use       usecases.Use
	templates *template.Template
}

func New(storage usecases.Storage) (*Cli, error) {
	templates, err := compileTemplates()
	if err != nil {
		return nil, err
	}

	return &Cli{use: usecases.New(storage), templates: templates}, nil
}

func MustNew(storage usecases.Storage) *Cli {
	cli, err := New(storage)
	if err != nil {
		panic(err)
	}

	return cli
}

func (cli *Cli) WithOverride(use usecases.Use) *Cli {
	cli.use = use

	return cli
}

func (cli *Cli) Dispatch(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.Usage()
	}

	command, args := args[0], args[1:]

	return cli.dispatch(ctx, command, args)
}

func (cli *Cli) dispatch(ctx context.Context, command string, args []string) int {
	switch command {
	case "add":
		return cli.Add(ctx, args)
	case "update":
		return cli.Update(ctx, args)
	case "delete":
		return cli.Delete(ctx, args)
	case "mark":
		return cli.Mark(ctx, args)
	case "work":
		return cli.Work(ctx, args)
	case "done":
		return cli.Done(ctx, args)
	case "list":
		return cli.List(ctx, args)
	case "help":
		return cli.Help()
	default:
		return cli.Unknown(command)
	}
}
