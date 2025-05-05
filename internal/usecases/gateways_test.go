package usecases_test

import (
	"context"
	"tasker/internal/domain"
	"tasker/pkg/testkit"
)

type StorageMock struct {
	Data map[uint64]*domain.Task

	SaveTaskFunc     func(ctx context.Context, task *domain.Task) (*domain.Task, error)
	UpdateTaskFunc   func(ctx context.Context, task *domain.Task) error
	DeleteTaskFunc   func(ctx context.Context, tid uint64) error
	GetByIDFunc      func(ctx context.Context, tid uint64) (*domain.Task, error)
	ListAllFunc      func(ctx context.Context) ([]*domain.Task, error)
	ListByStatusFunc func(ctx context.Context, status domain.Status) ([]*domain.Task, error)
}

func (s *StorageMock) SaveTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if s.SaveTaskFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.SaveTaskFunc(ctx, task)
}

func (s *StorageMock) UpdateTask(ctx context.Context, task *domain.Task) error {
	if s.UpdateTaskFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.UpdateTaskFunc(ctx, task)
}

func (s *StorageMock) DeleteTask(ctx context.Context, tid uint64) error {
	if s.DeleteTaskFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.DeleteTaskFunc(ctx, tid)
}

func (s *StorageMock) GetByID(ctx context.Context, tid uint64) (*domain.Task, error) {
	if s.GetByIDFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.GetByIDFunc(ctx, tid)
}

func (s *StorageMock) ListAll(ctx context.Context) ([]*domain.Task, error) {
	if s.ListAllFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.ListAllFunc(ctx)
}

func (s *StorageMock) ListByStatus(ctx context.Context, status domain.Status) ([]*domain.Task, error) {
	if s.ListByStatusFunc == nil {
		panic(testkit.ErrDummy)
	}

	return s.ListByStatusFunc(ctx, status)
}
