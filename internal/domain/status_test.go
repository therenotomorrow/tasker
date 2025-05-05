package domain_test

import (
	"errors"
	"testing"

	"github.com/therenotomorrow/tasker/internal/domain"
)

func TestUnitNewStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		raw string
	}

	type want struct {
		status domain.Status
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "todo", args: args{raw: "todo"}, want: want{status: domain.StatusTodo}},
		{name: "done", args: args{raw: "done"}, want: want{status: domain.StatusDone}},
		{name: "progress", args: args{raw: "progress"}, want: want{status: domain.StatusProgress}},
		{name: "invalid", args: args{raw: "invalid"}, want: want{err: domain.ErrInvalidStatus}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewStatus(test.args.raw)

			if !errors.Is(err, test.want.err) {
				t.Fatalf("NewStatus() error = %v, want = %v", err, test.want.err)
			}

			if got != test.want.status {
				t.Errorf("NewStatus() got = %v, want = %v", got, test.want.status)
			}
		})
	}
}
