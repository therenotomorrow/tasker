package cli

import (
	"context"
	"io"
	"os"
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

type Config struct {
	Output  io.Writer
	Storage usecases.Storage
}

type Cli struct {
	config    Config
	use       *usecases.UseCases
	templates *template.Template
}

func New(config Config) *Cli {
	use := usecases.New(config.Storage)
	templates := compileTemplates()

	if config.Output == nil {
		config.Output = os.Stdout
	}

	return &Cli{use: use, config: config, templates: templates}
}

func (cli *Cli) Dispatch(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.Usage()
	}

	command, args := args[0], args[1:]

	status := cli.dispatch(ctx, command, args)

	// end output with new line
	_, _ = cli.config.Output.Write([]byte{'\n'})

	return status
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
