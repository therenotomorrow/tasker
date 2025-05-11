package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/therenotomorrow/tasker/internal/cli"
)

const (
	// duplicate to be sure that external code will use exactly the same statuses.
	success = iota
	failure
	invalid
	unknown
	noArgs
)

func TestUnitNew(t *testing.T) {
	t.Parallel()

	type args struct {
		config cli.Config
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "nothing", args: args{config: cli.Config{}}},
		{name: "stdout", args: args{config: cli.Config{Output: os.Stdout}}},
		{name: "other", args: args{config: cli.Config{Output: new(strings.Builder)}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tacker := cli.New(test.args.config)

			if tacker == nil {
				t.Errorf("New() got = %v, want = %v", tacker, new(cli.Cli))
			}
		})
	}
}

func TestUnitCliDispatch(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "usage", args: args{args: make([]string, 0)}, want: noArgs},
		{name: "add", args: args{args: []string{"add"}}, want: noArgs},
		{name: "update", args: args{args: []string{"update"}}, want: noArgs},
		{name: "delete", args: args{args: []string{"delete"}}, want: noArgs},
		{name: "mark", args: args{args: []string{"mark"}}, want: noArgs},
		{name: "work", args: args{args: []string{"work"}}, want: noArgs},
		{name: "done", args: args{args: []string{"done"}}, want: noArgs},
		{name: "list", args: args{args: []string{"list", "invalid"}}, want: invalid},
		{name: "help", args: args{args: []string{"help"}}, want: success},
		{name: "unknown", args: args{args: []string{"unknown"}}, want: failure},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tasker := cli.New(cli.Config{})

			if got := tasker.Dispatch(ctx, test.args.args); got != test.want {
				t.Errorf("Dispatch() got = %v, want = %v", got, test.want)
			}
		})
	}

	ctx := t.Context()
	buffer := bytes.NewBuffer(nil)
	tasker := cli.New(cli.Config{Output: buffer})

	tasker.Dispatch(ctx, []string{"run"})

	// new line must be at the end of output
	want := "error: unknown command \"run\"\n"

	if got := buffer.String(); got != want {
		t.Errorf("Dispatch() got = %v, want = %v", got, want)
	}
}
