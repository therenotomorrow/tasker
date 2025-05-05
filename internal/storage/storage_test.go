package storage_test

import (
	"cmp"
	"errors"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/therenotomorrow/tasker/internal/config"
	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/internal/storage"
	"github.com/therenotomorrow/tasker/internal/usecases"
	"github.com/therenotomorrow/tasker/pkg/jsonfile"
	"github.com/therenotomorrow/tasker/pkg/testkit"
)

const ownerRO = 0o400

func TestUnitStorage(t *testing.T) {
	t.Parallel()

	var _ usecases.Storage = new(storage.Storage)
}

func TestIntegrationNew(t *testing.T) {
	t.Parallel()

	type args struct {
		config jsonfile.Config
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cannot find file",
			args: args{config: jsonfile.Config{File: string([]byte{0})}},
			want: "Storage error: JSONFile error: open",
		},
		{
			name: "corrupted file",
			args: args{config: jsonfile.Config{File: config.Path("test", "data", "empty.json")}},
			want: "Storage error: JSONFile decode error: EOF",
		},
		{
			name: testkit.SuccessTest,
			args: args{config: jsonfile.Config{File: config.Path("test", "data", "tasks-ro.json")}},
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stor, err := storage.New(test.args.config)

			if test.want != "" && !strings.Contains(err.Error(), test.want) {
				t.Fatalf("New() error = %v, want = %v", err, test.want)
			}

			if test.name != testkit.SuccessTest {
				return
			}

			got := stor.LastID()
			want := uint64(5)

			if got != want {
				t.Errorf("New() got = %v, want = %v", got, want)
			}
		})
	}
}

func TestIntegrationMustNew(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("MustNew() should panic")
		}
	}()

	_ = storage.MustNew(jsonfile.Config{File: "."})
}

func TestIntegrationStorageCRUDTask(t *testing.T) {
	t.Parallel()

	var (
		ctx      = t.Context()
		filename = config.Path("test", "data", "tasks-rw.json")
		stor     = storage.MustNew(jsonfile.Config{File: filename})
	)

	got, _ := stor.SaveTask(ctx, &domain.Task{ID: 666, Description: "sleeper", Status: domain.StatusTodo})
	want := &domain.Task{ID: 6, Description: "sleeper", Status: domain.StatusTodo}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SaveTask() got = %v, want = %v", got, want)
	}

	saved, _ := stor.GetByID(ctx, 6)
	if !reflect.DeepEqual(saved, want) {
		t.Errorf("SaveTask() got = %v, want = %v", saved, want)
	}

	list, _ := stor.ListAll(ctx)
	if want := 5; len(list) != want {
		t.Errorf("ListAll() got = %v, want = %v", len(list), want)
	}

	got.Description = "awaker"

	_ = stor.UpdateTask(ctx, got)
	saved, _ = stor.GetByID(ctx, 6)

	if !reflect.DeepEqual(saved, got) {
		t.Errorf("UpdateTask() got = %v, want = %v", saved, got)
	}

	_ = stor.DeleteTask(ctx, got)

	saved, _ = stor.GetByID(ctx, 6)
	if saved != nil {
		t.Errorf("DeleteTask() got = %v, want = %v", saved, got)
	}

	list, _ = stor.ListAll(ctx)
	if want := 4; len(list) != want {
		t.Errorf("ListAll() got = %v, want = %v", len(list), want)
	}
}

func TestIntegrationStorageListTasks(t *testing.T) {
	t.Parallel()

	var (
		ctx      = t.Context()
		filename = config.Path("test", "data", "tasks-ro.json")
		stor     = storage.MustNew(jsonfile.Config{File: filename})
		collect  = func(list []*domain.Task) []*domain.Task {
			ret := make([]*domain.Task, 0)

			for _, item := range list {
				ret = append(ret, &domain.Task{ID: item.ID, Description: item.Description, Status: item.Status})
			}

			slices.SortStableFunc(ret, func(a, b *domain.Task) int {
				return cmp.Compare(a.ID, b.ID)
			})

			return ret
		}
	)

	list, _ := stor.ListAll(ctx)
	if want := 4; len(list) != want {
		t.Errorf("ListAll() got = %v, want = %v", len(list), want)
	}

	list, _ = stor.ListByStatus(ctx, domain.StatusTodo)

	got := collect(list)
	want := []*domain.Task{{
		ID:          1,
		Description: "one",
		Status:      domain.StatusTodo,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListAll() got = %v, want = %v", got, want)
	}

	list, _ = stor.ListByStatus(ctx, domain.StatusProgress)

	got = collect(list)
	want = []*domain.Task{{
		ID:          3,
		Description: "three",
		Status:      domain.StatusProgress,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListAll() got = %v, want = %v", got, want)
	}

	list, _ = stor.ListByStatus(ctx, domain.StatusDone)

	got = collect(list)
	want = []*domain.Task{
		{ID: 4, Description: "four", Status: domain.StatusDone},
		{ID: 5, Description: "five", Status: domain.StatusDone},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListAll() got = %v, want = %v", got, want)
	}
}

func TestIntegrationStorageSaveTask(t *testing.T) {
	t.Parallel()

	const filename = "save-task.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
		task = &domain.Task{ID: 1, Status: domain.StatusTodo}
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	got, err := stor.SaveTask(ctx, task)
	if err == nil || got != nil {
		t.Fatalf("SaveTask() should failed, got = %v", err)
	}

	_ = os.Remove(filename)
	_ = os.WriteFile(filename, []byte("{}"), ownerRO)

	got, err = stor.SaveTask(ctx, task)
	if err == nil || got != nil {
		t.Fatalf("SaveTask() should failed, got = %v", err)
	}
}

func TestIntegrationStorageUpdateTask(t *testing.T) {
	t.Parallel()

	const filename = "update-task.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
		task = &domain.Task{ID: 1, Status: domain.StatusTodo}
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	err := stor.UpdateTask(ctx, task)
	if err == nil {
		t.Fatalf("UpdateTask() should failed, got = %v", err)
	}

	_ = os.Remove(filename)
	_ = os.WriteFile(filename, []byte("{}"), ownerRO)

	err = stor.UpdateTask(ctx, task)
	if err == nil {
		t.Fatalf("UpdateTask() should failed, got = %v", err)
	}
}

func TestIntegrationStorageDeleteTask(t *testing.T) {
	t.Parallel()

	const filename = "delete-task.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
		task = &domain.Task{ID: 1, Status: domain.StatusTodo}
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	err := stor.DeleteTask(ctx, task)
	if err == nil {
		t.Fatalf("DeleteTask() should failed, got = %v", err)
	}

	_ = os.Remove(filename)
	_ = os.WriteFile(filename, []byte("{}"), ownerRO)

	err = stor.DeleteTask(ctx, task)
	if err == nil {
		t.Fatalf("DeleteTask() should failed, got = %v", err)
	}
}

func TestIntegrationStorageGetByID(t *testing.T) {
	t.Parallel()

	const filename = "get-by-id.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
		tid  = uint64(1)
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	got, err := stor.GetByID(ctx, tid)
	if err == nil || got != nil {
		t.Fatalf("GetByID() should failed, got = %v", err)
	}

	_ = os.Remove(filename)
	_ = os.WriteFile(filename, []byte("{}"), 0o600)

	got, err = stor.GetByID(ctx, tid)
	if !errors.Is(err, domain.ErrTaskNotFound) || got != nil {
		t.Fatalf("GetByID() error = %v, want = %v", err, domain.ErrTaskNotFound)
	}
}

func TestIntegrationStorageListAll(t *testing.T) {
	t.Parallel()

	const filename = "list-all.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	got, err := stor.ListAll(ctx)
	if err == nil || len(got) != 0 {
		t.Fatalf("ListAll() should failed, got = %v", err)
	}
}

func TestIntegrationStorageListByStatus(t *testing.T) {
	t.Parallel()

	const filename = "list-by-status.json"

	var (
		ctx  = t.Context()
		stor = storage.MustNew(jsonfile.Config{File: filename})
	)

	defer func() { _ = os.Remove(filename) }()

	_, _ = os.Create(filename)

	got, err := stor.ListByStatus(ctx, domain.StatusDone)
	if err == nil || len(got) != 0 {
		t.Fatalf("ListByStatus() should failed, got = %v", err)
	}
}
