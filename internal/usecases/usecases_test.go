package usecases_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/internal/storage"
	"github.com/therenotomorrow/tasker/internal/usecases"
	"github.com/therenotomorrow/tasker/pkg/testkit"
)

const taskNotFoundTest = "task not found"

func TestUnitNew(t *testing.T) {
	t.Parallel()

	use := usecases.New(new(storage.Mock))

	if use == nil {
		t.Errorf("New() got = %v, want = %v", use, new(usecases.UseCases))
	}
}

func TestUnitUseCasesAddTask(t *testing.T) {
	t.Parallel()

	type args struct {
		description string
	}

	type want struct {
		task *domain.Task
		err  error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty description",
			args: args{description: ""},
			want: want{err: domain.ErrEmptyDescription},
		},
		{
			name: "empty trimmed description",
			args: args{description: "    "},
			want: want{err: domain.ErrEmptyDescription},
		},
		{
			name: testkit.FailureTest,
			args: args{description: "some description"},
			want: want{err: testkit.ErrDummy},
		},
		{
			name: testkit.SuccessTest,
			args: args{description: "  some task here  "},
			want: want{task: &domain.Task{
				ID:          1,
				Description: "some task here",
				Status:      domain.StatusTodo,
				CreatedAt:   time.Now().Truncate(time.Minute),
				UpdatedAt:   time.Now().Truncate(time.Minute),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor := new(storage.Mock)

			stor.SaveTaskFunc = func(ctx context.Context, task *domain.Task) (*domain.Task, error) {
				var err error

				switch test.name {
				case testkit.SuccessTest:
					task.ID = 1
				case testkit.FailureTest:
					err = testkit.ErrDummy
				}

				return task, err
			}

			ctx := t.Context()
			use := usecases.New(stor)
			got, err := use.AddTask(ctx, test.args.description)

			if !errors.Is(err, test.want.err) {
				t.Fatalf("AddTask() error = %v, want = %v", err, test.want.err)
			}

			if got != nil {
				got.CreatedAt = got.CreatedAt.Truncate(time.Minute)
				got.UpdatedAt = got.UpdatedAt.Truncate(time.Minute)
			}

			if !reflect.DeepEqual(got, test.want.task) {
				t.Errorf("AddTask() got = %v, want = %v", got, test.want.task)
			}
		})
	}
}

func TestUnitUseCasesUpdateTask(t *testing.T) {
	t.Parallel()

	type args struct {
		tid         string
		description string
	}

	type want struct {
		task *domain.Task
		err  error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid taskID",
			args: args{tid: "invalid", description: ""},
			want: want{err: domain.ErrInvalidTaskID},
		},
		{
			name: "negative taskID",
			args: args{tid: "-1", description: ""},
			want: want{err: domain.ErrInvalidTaskID},
		},
		{
			name: "empty description",
			args: args{tid: "1", description: ""},
			want: want{err: domain.ErrEmptyDescription},
		},
		{
			name: "empty trimmed description",
			args: args{tid: "1", description: "    "},
			want: want{err: domain.ErrEmptyDescription},
		},
		{
			name: taskNotFoundTest,
			args: args{tid: "0", description: "some description"},
			want: want{err: domain.ErrTaskNotFound},
		},
		{
			name: testkit.FailureTest,
			args: args{tid: "1", description: "some description"},
			want: want{err: testkit.ErrDummy},
		},
		{
			name: testkit.SuccessTest,
			args: args{tid: "1", description: "  new description  "},
			want: want{task: &domain.Task{
				ID:          1,
				Description: "new description",
				Status:      domain.StatusDone,
				UpdatedAt:   time.Now().Truncate(time.Minute),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor := new(storage.Mock)

			stor.GetByIDFunc = func(ctx context.Context, tid uint64) (*domain.Task, error) {
				if test.name == taskNotFoundTest {
					return nil, domain.ErrTaskNotFound
				}

				return &domain.Task{ID: 1, Description: "old description", Status: domain.StatusDone}, nil
			}
			stor.UpdateTaskFunc = func(ctx context.Context, task *domain.Task) error {
				if test.name == testkit.FailureTest {
					return testkit.ErrDummy
				}

				return nil
			}

			ctx := t.Context()
			use := usecases.New(stor)
			got, err := use.UpdateTask(ctx, test.args.tid, test.args.description)

			if !errors.Is(err, test.want.err) {
				t.Fatalf("UpdateTask() error = %v, want = %v", err, test.want.err)
			}

			if got != nil {
				got.UpdatedAt = got.UpdatedAt.Truncate(time.Minute)
			}

			if !reflect.DeepEqual(got, test.want.task) {
				t.Errorf("UpdateTask() got = %v, want = %v", got, test.want.task)
			}
		})
	}
}

func TestUnitUseCasesDeleteTask(t *testing.T) {
	t.Parallel()

	type args struct {
		tid string
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{name: "invalid taskID", args: args{tid: "invalid"}, want: domain.ErrInvalidTaskID},
		{name: "negative taskID", args: args{tid: "-1"}, want: domain.ErrInvalidTaskID},
		{name: taskNotFoundTest, args: args{tid: "0"}, want: domain.ErrTaskNotFound},
		{name: testkit.FailureTest, args: args{tid: "1"}, want: testkit.ErrDummy},
		{name: testkit.SuccessTest, args: args{tid: "1"}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor := new(storage.Mock)

			stor.GetByIDFunc = func(ctx context.Context, tid uint64) (*domain.Task, error) {
				if test.name == taskNotFoundTest {
					return nil, domain.ErrTaskNotFound
				}

				return &domain.Task{ID: 1, Description: "description", Status: domain.StatusDone}, nil
			}
			stor.DeleteTaskFunc = func(ctx context.Context, task *domain.Task) error {
				if test.name == testkit.FailureTest {
					return testkit.ErrDummy
				}

				return nil
			}

			ctx := t.Context()
			use := usecases.New(stor)
			err := use.DeleteTask(ctx, test.args.tid)

			if !errors.Is(err, test.want) {
				t.Fatalf("DeleteTask() error = %v, want = %v", err, test.want)
			}
		})
	}
}

func TestUnitUseCasesMarkTask(t *testing.T) {
	t.Parallel()

	type args struct {
		tid  string
		mark string
	}

	type want struct {
		task *domain.Task
		err  error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid taskID",
			args: args{tid: "invalid", mark: ""},
			want: want{err: domain.ErrInvalidTaskID},
		},
		{
			name: "negative taskID",
			args: args{tid: "-1", mark: ""},
			want: want{err: domain.ErrInvalidTaskID},
		},
		{
			name: "empty status",
			args: args{tid: "1", mark: ""},
			want: want{err: domain.ErrInvalidStatus},
		},
		{
			name: "invalid status",
			args: args{tid: "1", mark: "invalid"},
			want: want{err: domain.ErrInvalidStatus},
		},
		{
			name: taskNotFoundTest,
			args: args{tid: "0", mark: "todo"},
			want: want{err: domain.ErrTaskNotFound},
		},
		{
			name: "task is done",
			args: args{tid: "1", mark: "progress"},
			want: want{err: domain.ErrTaskAlreadyDone},
		},
		{
			name: testkit.FailureTest,
			args: args{tid: "1", mark: "progress"},
			want: want{err: testkit.ErrDummy},
		},
		{
			name: testkit.SuccessTest,
			args: args{tid: "1", mark: "done"},
			want: want{task: &domain.Task{
				ID:        1,
				Status:    domain.StatusDone,
				UpdatedAt: time.Now().Truncate(time.Minute),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor := new(storage.Mock)

			stor.GetByIDFunc = func(ctx context.Context, tid uint64) (*domain.Task, error) {
				switch test.name {
				case taskNotFoundTest:
					return nil, domain.ErrTaskNotFound
				case "task is done":
					return &domain.Task{ID: 1, Status: domain.StatusDone}, nil
				}

				return &domain.Task{ID: 1, Status: domain.StatusTodo}, nil
			}
			stor.UpdateTaskFunc = func(ctx context.Context, task *domain.Task) error {
				if test.name == testkit.FailureTest {
					return testkit.ErrDummy
				}

				return nil
			}

			ctx := t.Context()
			use := usecases.New(stor)
			got, err := use.MarkTask(ctx, test.args.tid, test.args.mark)

			if !errors.Is(err, test.want.err) {
				t.Fatalf("MarkTask() error = %v, want = %v", err, test.want.err)
			}

			if got != nil {
				got.UpdatedAt = got.UpdatedAt.Truncate(time.Minute)
			}

			if !reflect.DeepEqual(got, test.want.task) {
				t.Errorf("MarkTask() got = %v, want = %v", got, test.want.task)
			}
		})
	}
}

func TestUnitUseCasesListTasks(t *testing.T) {
	t.Parallel()

	type args struct {
		params usecases.ListParams
	}

	type want struct {
		tasks []*domain.Task
		err   error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid status",
			args: args{params: usecases.ListParams{Status: "invalid"}},
			want: want{err: domain.ErrInvalidStatus},
		},
		{
			name: "list all success",
			args: args{params: usecases.ListParams{Status: ""}},
			want: want{tasks: []*domain.Task{
				{ID: 1, Status: domain.StatusDone},
				{ID: 2, Status: domain.StatusTodo},
			}},
		},
		{
			name: "list all failure",
			args: args{params: usecases.ListParams{Status: ""}},
			want: want{err: testkit.ErrDummy},
		},
		{
			name: "list by status success",
			args: args{params: usecases.ListParams{Status: "todo"}},
			want: want{tasks: []*domain.Task{
				{ID: 1, Status: domain.StatusTodo},
				{ID: 2, Status: domain.StatusTodo},
			}},
		},
		{
			name: "list by status failure",
			args: args{params: usecases.ListParams{Status: "done"}},
			want: want{err: testkit.ErrDummy},
		},
		{
			name: "empty list",
			args: args{params: usecases.ListParams{Status: ""}},
			want: want{err: domain.ErrEmptyTasks},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor := new(storage.Mock)

			stor.ListAllFunc = func(ctx context.Context) ([]*domain.Task, error) {
				switch test.name {
				case "empty list":
					return make([]*domain.Task, 0), nil
				case "list all failure":
					return nil, testkit.ErrDummy
				}

				return []*domain.Task{{ID: 1, Status: domain.StatusDone}, {ID: 2, Status: domain.StatusTodo}}, nil
			}
			stor.ListByStatusFunc = func(ctx context.Context, status domain.Status) ([]*domain.Task, error) {
				if test.name == "list by status failure" {
					return nil, testkit.ErrDummy
				}

				if status != domain.StatusTodo {
					panic("invalid status for mock")
				}

				return []*domain.Task{
					{ID: 1, Status: domain.StatusTodo},
					{ID: 2, Status: domain.StatusTodo},
				}, nil
			}

			ctx := t.Context()
			use := usecases.New(stor)
			got, err := use.ListTasks(ctx, test.args.params)

			if !errors.Is(err, test.want.err) {
				t.Fatalf("ListTasks() error = %v, want = %v", err, test.want.err)
			}

			if !reflect.DeepEqual(got, test.want.tasks) {
				t.Errorf("ListTasks() got = %v, want = %v", got, test.want.tasks)
			}
		})
	}
}
