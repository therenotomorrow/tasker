package cli_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/therenotomorrow/tasker/internal/cli"
	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/internal/storage"
	"github.com/therenotomorrow/tasker/internal/usecases"
	"github.com/therenotomorrow/tasker/pkg/testkit"
)

const taskNotFoundTest = "task not found"

func newMock(testName string) *storage.Mock {
	stor := new(storage.Mock)
	stor.SaveTaskFunc = func(ctx context.Context, task *domain.Task) (*domain.Task, error) {
		if testName == testkit.SuccessTest {
			return task, nil
		}

		return nil, testkit.ErrDummy
	}
	stor.GetByIDFunc = func(ctx context.Context, tid uint64) (*domain.Task, error) {
		switch testName {
		case taskNotFoundTest:
			return nil, domain.ErrTaskNotFound
		case "task is done":
			return &domain.Task{Status: domain.StatusDone}, nil
		}

		return new(domain.Task), nil
	}
	stor.UpdateTaskFunc = func(ctx context.Context, task *domain.Task) error {
		if testName == testkit.SuccessTest {
			return nil
		}

		return testkit.ErrDummy
	}
	stor.DeleteTaskFunc = func(ctx context.Context, task *domain.Task) error {
		if testName == testkit.SuccessTest {
			return nil
		}

		return testkit.ErrDummy
	}
	stor.ListAllFunc = func(ctx context.Context) ([]*domain.Task, error) {
		switch testName {
		case "empty list":
			return make([]*domain.Task, 0), nil
		case testkit.SuccessTest:
			return []*domain.Task{
				{
					ID:        2,
					Status:    domain.StatusDone,
					CreatedAt: time.Date(1992, 1, 28, 15, 45, 12, 0, time.UTC),
					UpdatedAt: time.Now().Add(-30 * time.Minute),
				},
				{
					ID:        3,
					Status:    domain.StatusProgress,
					CreatedAt: time.Date(1998, 4, 20, 1, 43, 12, 0, time.UTC),
					UpdatedAt: time.Now().Add(-2 * time.Hour),
				},
				{
					ID:        1,
					Status:    domain.StatusTodo,
					CreatedAt: time.Date(1970, 4, 20, 0, 2, 12, 0, time.UTC),
					UpdatedAt: time.Now(),
				},
				{
					ID:        4,
					Status:    domain.StatusDone,
					CreatedAt: time.Date(2003, 5, 8, 10, 10, 10, 0, time.UTC),
					UpdatedAt: time.Now().Add(-25 * time.Hour),
				},
			}, nil
		}

		return nil, testkit.ErrDummy
	}
	stor.ListByStatusFunc = func(ctx context.Context, status domain.Status) ([]*domain.Task, error) {
		return []*domain.Task{{
			ID:        1,
			Status:    domain.StatusTodo,
			CreatedAt: time.Date(1992, 5, 8, 10, 10, 10, 0, time.UTC),
			UpdatedAt: time.Now().Add(-time.Hour),
		}}, nil
	}

	return stor
}

func newCli(out io.Writer, mock usecases.Storage) *cli.Cli {
	return cli.New(cli.Config{Output: out, Storage: mock})
}

func TestUnitCliAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		text string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not enough arguments",
			args: args{args: make([]string, 0)},
			want: want{code: noArgs, text: `error: not enough arguments for command "add"`},
		},
		{
			name: "invalid description",
			args: args{args: []string{"    "}},
			want: want{code: invalid, text: `error: invalid "description" parameter, must be not empty`},
		},
		{
			name: "unexpected error",
			args: args{args: []string{"description"}},
			want: want{code: unknown, text: `error: unexpected behaviour "AddTask error: dummy"`},
		},
		{
			name: testkit.SuccessTest,
			args: args{args: []string{"description"}},
			want: want{code: success, text: `task added successfully (ID: 0)`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			got := client.Add(ctx, test.args.args)

			if got != test.want.code {
				t.Errorf("Add() got = %v, want = %v", got, test.want.code)
			}

			if text := buffer.String(); text != test.want.text {
				t.Errorf("Add() got = %v, want = %v", text, test.want.text)
			}
		})
	}
}

func TestUnitCliUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		text string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not enough arguments",
			args: args{args: make([]string, 0)},
			want: want{code: noArgs, text: `error: not enough arguments for command "update"`},
		},
		{
			name: "invalid description",
			args: args{args: []string{"1", "    "}},
			want: want{code: invalid, text: `error: invalid "description" parameter, must be not empty`},
		},
		{
			name: "invalid task ID",
			args: args{args: []string{"one", "description"}},
			want: want{code: invalid, text: `error: invalid "id" parameter, must be positive integer`},
		},
		{
			name: taskNotFoundTest,
			args: args{args: []string{"1", "description"}},
			want: want{code: failure, text: `error: task (ID: 1) not found`},
		},
		{
			name: "unexpected error",
			args: args{args: []string{"1", "description"}},
			want: want{code: unknown, text: `error: unexpected behaviour "UpdateTask error: dummy"`},
		},
		{
			name: testkit.SuccessTest,
			args: args{args: []string{"1", "description"}},
			want: want{code: success, text: `task updated successfully`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			got := client.Update(ctx, test.args.args)

			if got != test.want.code {
				t.Errorf("Update() got = %v, want = %v", got, test.want.code)
			}

			if text := buffer.String(); text != test.want.text {
				t.Errorf("Update() got = %v, want = %v", text, test.want.text)
			}
		})
	}
}

func TestUnitCliDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		text string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not enough arguments",
			args: args{args: make([]string, 0)},
			want: want{code: noArgs, text: `error: not enough arguments for command "delete"`},
		},
		{
			name: "invalid task ID",
			args: args{args: []string{"one"}},
			want: want{code: invalid, text: `error: invalid "id" parameter, must be positive integer`},
		},
		{
			name: taskNotFoundTest,
			args: args{args: []string{"1"}},
			want: want{code: failure, text: `error: task (ID: 1) not found`},
		},
		{
			name: "unexpected error",
			args: args{args: []string{"1"}},
			want: want{code: unknown, text: `error: unexpected behaviour "DeleteTask error: dummy"`},
		},
		{
			name: testkit.SuccessTest,
			args: args{args: []string{"1"}},
			want: want{code: success, text: `task deleted successfully`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			got := client.Delete(ctx, test.args.args)

			if got != test.want.code {
				t.Errorf("Delete() got = %v, want = %v", got, test.want.code)
			}

			if text := buffer.String(); text != test.want.text {
				t.Errorf("Delete() got = %v, want = %v", text, test.want.text)
			}
		})
	}
}

func TestUnitCliMark(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		text string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not enough arguments",
			args: args{args: make([]string, 0)},
			want: want{code: noArgs, text: `error: not enough arguments for command "mark"`},
		},
		{
			name: "invalid task ID",
			args: args{args: []string{"one", "todo"}},
			want: want{code: invalid, text: `error: invalid "id" parameter, must be positive integer`},
		},
		{
			name: "invalid status",
			args: args{args: []string{"1", "invalid"}},
			want: want{code: invalid, text: `error: invalid "status" parameter, must be one of [todo progress done]`},
		},
		{
			name: taskNotFoundTest,
			args: args{args: []string{"1", "todo"}},
			want: want{code: failure, text: `error: task (ID: 1) not found`},
		},
		{
			name: "task is done",
			args: args{args: []string{"1", "todo"}},
			want: want{code: failure, text: `error: cannot change status for done task`},
		},
		{
			name: "unexpected error",
			args: args{args: []string{"1", "todo"}},
			want: want{code: unknown, text: `error: unexpected behaviour "MarkTask error: dummy"`},
		},
		{
			name: testkit.SuccessTest,
			args: args{args: []string{"1", "todo"}},
			want: want{code: success, text: `task status changed successfully`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			got := client.Mark(ctx, test.args.args)

			if got != test.want.code {
				t.Errorf("Mark() got = %v, want = %v", got, test.want.code)
			}

			if text := buffer.String(); text != test.want.text {
				t.Errorf("Mark() got = %v, want = %v", text, test.want.text)
			}
		})
	}
}

func TestUnitCliWorkDoneShortcuts(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		work string
		done string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not enough arguments",
			args: args{args: make([]string, 0)},
			want: want{
				code: noArgs,
				work: `error: not enough arguments for command "work"`,
				done: `error: not enough arguments for command "done"`,
			},
		},
		{
			name: "invalid task ID",
			args: args{args: []string{"one", "todo"}},
			want: want{
				code: invalid,
				work: `error: invalid "id" parameter, must be positive integer`,
				done: `error: invalid "id" parameter, must be positive integer`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			done := client.Done(ctx, test.args.args)

			if done != test.want.code {
				t.Errorf("Done() got = %v, want = %v", done, test.want.code)
			}

			if text := buffer.String(); text != test.want.done {
				t.Errorf("Done() got = %v, want = %v", text, test.want.done)
			}

			buffer.Reset()

			work := client.Work(ctx, test.args.args)

			if work != test.want.code {
				t.Errorf("Work() got = %v, want = %v", work, test.want.code)
			}

			if text := buffer.String(); text != test.want.work {
				t.Errorf("Work() got = %v, want = %v", text, test.want.work)
			}
		})
	}
}

func TestUnitCliList(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}

	type want struct {
		code int
		text string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid status",
			args: args{args: []string{"invalid"}},
			want: want{code: invalid, text: `error: invalid "status" parameter, must be one of [todo progress done]`},
		},
		{
			name: "empty list",
			args: args{args: make([]string, 0)},
			want: want{code: failure, text: `error: task list is empty`},
		},
		{
			name: "unexpected error",
			args: args{args: make([]string, 0)},
			want: want{code: unknown, text: `error: unexpected behaviour "ListTasks error: dummy"`},
		},
		{
			name: testkit.SuccessTest,
			args: args{args: make([]string, 0)},
			want: want{code: success, text: `
---- id: 1
description | 
status      | todo
created at  | 20 Apr 1970 00:02:12
last update | 0 minute(s) ago

---- id: 2
description | 
status      | done
created at  | 28 Jan 1992 15:45:12
last update | 30 minute(s) ago

---- id: 3
description | 
status      | progress
created at  | 20 Apr 1998 01:43:12
last update | 2 hour(s) ago

---- id: 4
description | 
status      | done
created at  | 08 May 2003 10:10:10
last update | 1 day(s) ago
`},
		},
		{
			name: "with status",
			args: args{args: []string{"todo"}},
			want: want{code: success, text: `
---- id: 1
description | 
status      | todo
created at  | 08 May 1992 10:10:10
last update | 1 hour(s) ago
`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			buffer := bytes.NewBuffer(nil)
			client := newCli(buffer, newMock(test.name))

			got := client.List(ctx, test.args.args)

			if got != test.want.code {
				t.Errorf("List() got = %v, want = %v", got, test.want.code)
			}

			if text := buffer.String(); text != test.want.text {
				t.Errorf("List() got = %v, want = %v", text, test.want.text)
			}
		})
	}
}

func TestUnitCliHelp(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)
	client := newCli(buffer, newMock(t.Name()))

	got := client.Help()

	if want := success; got != want {
		t.Errorf("Help() got = %v, want = %v", got, want)
	}

	want := `manage tasks with ease from the command line:
 - tasker add "description"
      add a new task with the given description
 - tasker update <id> "new description"
      update the description of an existing task by its ID
 - tasker delete <id>
      delete the task with the specified ID
 - tasker mark <id> <status>
      set a new status for the task ("todo", "progress", or "done")
 - tasker work <id>
      shortcut to mark the task as "progress"
 - tasker done <id>
      shortcut to mark the task as "done"
 - tasker list [status]
      list all tasks, if a status is provided, only tasks with that status will be shown
 - tasker help
      show this help message and exit`

	if text := buffer.String(); text != want {
		t.Errorf("Help() got = %v, want = %v", text, want)
	}
}

func TestUnitCliUsage(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)
	client := newCli(buffer, newMock(t.Name()))

	got := client.Usage()

	if want := noArgs; got != want {
		t.Errorf("Usage() got = %v, want = %v", got, want)
	}

	want := `manage tasks with ease from the command line:
 - tasker add "description"
      add a new task with the given description
 - tasker update <id> "new description"
      update the description of an existing task by its ID
 - tasker delete <id>
      delete the task with the specified ID
 - tasker mark <id> <status>
      set a new status for the task ("todo", "progress", or "done")
 - tasker work <id>
      shortcut to mark the task as "progress"
 - tasker done <id>
      shortcut to mark the task as "done"
 - tasker list [status]
      list all tasks, if a status is provided, only tasks with that status will be shown
 - tasker help
      show this help message and exit`

	if text := buffer.String(); text != want {
		t.Errorf("Usage() got = %v, want = %v", text, want)
	}
}

func TestUnitCliUnknown(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)
	client := newCli(buffer, newMock(t.Name()))

	got := client.Unknown("stop")

	if want := failure; got != want {
		t.Errorf("Unknown() got = %v, want = %v", got, want)
	}

	want := `error: unknown command "stop"`

	if text := buffer.String(); text != want {
		t.Errorf("Unknown() got = %v, want = %v", text, want)
	}
}
