package usecases_test

import (
	"context"
	"errors"
	"reflect"
	"tasker/internal/domain"
	"tasker/internal/usecases"
	"tasker/pkg/testkit"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	use := usecases.New(new(StorageMock))

	if use == nil {
		t.Errorf("New() got = %v, want = %v", use, new(usecases.UseCases))
	}
}

func TestUseCasesAddTask(t *testing.T) {
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
			want: want{task: nil, err: domain.ErrEmptyDescription},
		},
		{
			name: "empty trimmed description",
			args: args{description: "    "},
			want: want{task: nil, err: domain.ErrEmptyDescription},
		},
		{
			name: "failure",
			args: args{description: "some description"},
			want: want{task: nil, err: testkit.ErrDummy},
		},
		{
			name: "success",
			args: args{description: "  some task here  "},
			want: want{task: &domain.Task{
				ID:          1,
				Description: "some task here",
				Status:      domain.StatusTodo,
				CreatedAt:   time.Now().Truncate(time.Minute),
				UpdatedAt:   time.Now().Truncate(time.Minute),
			}, err: nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			storage := new(StorageMock)

			storage.SaveTaskFunc = func(ctx context.Context, task *domain.Task) (*domain.Task, error) {
				var err error

				switch test.name {
				case "success":
					task.ID = 1
				case "failure":
					err = testkit.ErrDummy
				}

				return task, err
			}

			ctx := t.Context()
			use := usecases.New(storage)
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
