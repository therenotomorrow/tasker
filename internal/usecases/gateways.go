package usecases

import (
	"context"

	"github.com/therenotomorrow/tasker/internal/domain"
)

type Saver interface {
	SaveTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type Updater interface {
	UpdateTask(ctx context.Context, task *domain.Task) error
}

type Deleter interface {
	DeleteTask(ctx context.Context, task *domain.Task) error
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

type Use interface {
	AddTask(ctx context.Context, description string) (*domain.Task, error)
	UpdateTask(ctx context.Context, tid string, description string) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
	MarkTask(ctx context.Context, tid string, mark string) (*domain.Task, error)
	ListTasks(ctx context.Context, params ListParams) ([]*domain.Task, error)
}
