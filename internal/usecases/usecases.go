package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/therenotomorrow/tasker/internal/domain"
)

type UseCases struct {
	storage Storage
}

func New(storage Storage) *UseCases {
	return &UseCases{storage: storage}
}

func (use *UseCases) AddTask(ctx context.Context, description string) (*domain.Task, error) {
	const where = "AddTask"

	description, err := use.validateDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	now := time.Now()
	task := &domain.Task{
		ID:          0,
		Description: description,
		Status:      domain.StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task, err = use.storage.SaveTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	return task, nil
}

func (use *UseCases) UpdateTask(ctx context.Context, tid string, description string) (*domain.Task, error) {
	const where = "UpdateTask"

	taskID, err := use.validateTaskID(tid)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	description, err = use.validateDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	task, err := use.storage.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	task.Description = description
	task.UpdatedAt = time.Now()

	err = use.storage.UpdateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	return task, nil
}

func (use *UseCases) DeleteTask(ctx context.Context, tid string) error {
	const where = "DeleteTask"

	taskID, err := use.validateTaskID(tid)
	if err != nil {
		return fmt.Errorf("%s error: %w", where, err)
	}

	task, err := use.storage.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("%s error: %w", where, err)
	}

	err = use.storage.DeleteTask(ctx, task)
	if err != nil {
		return fmt.Errorf("%s error: %w", where, err)
	}

	return nil
}

func (use *UseCases) MarkTask(ctx context.Context, tid string, mark string) (*domain.Task, error) {
	const where = "MarkTask"

	taskID, err := use.validateTaskID(tid)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	status, err := use.validateStatus(mark)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	task, err := use.storage.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	if task.IsDone() {
		return nil, fmt.Errorf("%s error: %w", where, domain.ErrTaskAlreadyDone)
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	err = use.storage.UpdateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	return task, nil
}

type ListParams struct {
	Status string
}

func (use *UseCases) ListTasks(ctx context.Context, params ListParams) ([]*domain.Task, error) {
	const where = "ListTasks"

	var err error

	status := domain.Status("")
	if params.Status != "" {
		status, err = use.validateStatus(params.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	var tasks []*domain.Task

	if status == "" {
		tasks, err = use.storage.ListAll(ctx)
	} else {
		tasks, err = use.storage.ListByStatus(ctx, status)
	}

	if err != nil {
		return nil, fmt.Errorf("%s error: %w", where, err)
	}

	if len(tasks) == 0 {
		return nil, domain.ErrEmptyTasks
	}

	return tasks, nil
}
