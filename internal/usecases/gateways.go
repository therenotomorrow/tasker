package usecases

import (
	"context"
	"tasker/internal/domain"
)

type Saver interface {
	SaveTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type Updater interface {
	UpdateTask(ctx context.Context, task *domain.Task) error
}

type Deleter interface {
	DeleteTask(ctx context.Context, tid uint64) error
}

type Retriever interface {
	GetByID(ctx context.Context, tid uint64) (*domain.Task, error)
	ListAll(ctx context.Context) ([]*domain.Task, error)
	ListByStatus(ctx context.Context, status domain.Status) ([]*domain.Task, error)
}

type Storage interface {
	Saver
	Updater
	Deleter
	Retriever
}
