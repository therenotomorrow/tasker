package storage_test

import (
	"context"
	"testing"

	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/internal/storage"
)

func TestUnitMockSaveTask(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.SaveTaskFunc = func(ctx context.Context, task *domain.Task) (*domain.Task, error) {
		return task, nil
	}

	_, _ = stor.SaveTask(ctx, new(domain.Task))

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("SaveTask() should panic")
		}
	}()

	stor.SaveTaskFunc = nil
	_, _ = stor.SaveTask(ctx, new(domain.Task))
}

func TestUnitMockUpdateTask(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.UpdateTaskFunc = func(ctx context.Context, task *domain.Task) error {
		return nil
	}

	_ = stor.UpdateTask(ctx, new(domain.Task))

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("UpdateTask() should panic")
		}
	}()

	stor.UpdateTaskFunc = nil
	_ = stor.UpdateTask(ctx, new(domain.Task))
}

func TestUnitMockDeleteTask(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.DeleteTaskFunc = func(ctx context.Context, task *domain.Task) error {
		return nil
	}

	_ = stor.DeleteTask(ctx, new(domain.Task))

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("DeleteTask() should panic")
		}
	}()

	stor.DeleteTaskFunc = nil
	_ = stor.DeleteTask(ctx, new(domain.Task))
}

func TestUnitMockGetByID(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.GetByIDFunc = func(ctx context.Context, tid uint64) (*domain.Task, error) {
		return new(domain.Task), nil
	}

	_, _ = stor.GetByID(ctx, 1)

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("GetByID() should panic")
		}
	}()

	stor.GetByIDFunc = nil
	_, _ = stor.GetByID(ctx, 1)
}

func TestUnitMockListAll(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.ListAllFunc = func(ctx context.Context) ([]*domain.Task, error) {
		return nil, nil
	}

	_, _ = stor.ListAll(ctx)

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("ListAll() should panic")
		}
	}()

	stor.ListAllFunc = nil
	_, _ = stor.ListAll(ctx)
}

func TestUnitMockListByStatus(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	stor := new(storage.Mock)
	stor.ListByStatusFunc = func(ctx context.Context, status domain.Status) ([]*domain.Task, error) {
		return nil, nil
	}

	_, _ = stor.ListByStatus(ctx, domain.StatusTodo)

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("ListByStatus() should panic")
		}
	}()

	stor.ListByStatusFunc = nil
	_, _ = stor.ListByStatus(ctx, domain.StatusTodo)
}
