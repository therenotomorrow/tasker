package storage

import (
	"context"

	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/pkg/testkit"
)

type Mock struct {
	SaveTaskFunc     func(ctx context.Context, task *domain.Task) (*domain.Task, error)
	UpdateTaskFunc   func(ctx context.Context, task *domain.Task) error
	DeleteTaskFunc   func(ctx context.Context, task *domain.Task) error
	GetByIDFunc      func(ctx context.Context, tid uint64) (*domain.Task, error)
	ListAllFunc      func(ctx context.Context) ([]*domain.Task, error)
	ListByStatusFunc func(ctx context.Context, status domain.Status) ([]*domain.Task, error)
}

func (s *Mock) SaveTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if s.SaveTaskFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.SaveTaskFunc(ctx, task)
}

func (s *Mock) UpdateTask(ctx context.Context, task *domain.Task) error {
	if s.UpdateTaskFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.UpdateTaskFunc(ctx, task)
}

func (s *Mock) DeleteTask(ctx context.Context, task *domain.Task) error {
	if s.DeleteTaskFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.DeleteTaskFunc(ctx, task)
}

func (s *Mock) GetByID(ctx context.Context, tid uint64) (*domain.Task, error) {
	if s.GetByIDFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.GetByIDFunc(ctx, tid)
}

func (s *Mock) ListAll(ctx context.Context) ([]*domain.Task, error) {
	if s.ListAllFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.ListAllFunc(ctx)
}

func (s *Mock) ListByStatus(ctx context.Context, status domain.Status) ([]*domain.Task, error) {
	if s.ListByStatusFunc == nil {
		panic(testkit.ErrUnimplemented)
	}

	return s.ListByStatusFunc(ctx, status)
}
