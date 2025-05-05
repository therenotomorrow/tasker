package domain_test

import (
	"testing"

	"github.com/therenotomorrow/tasker/internal/domain"
)

func TestUnitTaskIsDone(t *testing.T) {
	t.Parallel()

	type fields struct {
		Status domain.Status
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "done", fields: fields{Status: domain.StatusDone}, want: true},
		{name: "todo", fields: fields{Status: domain.StatusTodo}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := &domain.Task{Status: test.fields.Status}

			if got := task.IsDone(); got != test.want {
				t.Errorf("IsDone() got = %v, want = %v", got, test.want)
			}
		})
	}
}
