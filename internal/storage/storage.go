package storage

import (
	"context"
	"fmt"
	"tasker/internal/domain"
	"tasker/pkg/jsonfile"
)

const name = "Storage"

type Storage struct {
	fs     *jsonfile.JSONFile[Tasks]
	lastID uint64
}

func New(config jsonfile.Config) (*Storage, error) {
	jsonfs, err := jsonfile.New[Tasks](config)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	tasks, err := jsonfs.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	var lastID uint64

	for _, task := range tasks {
		lastID = max(lastID, task.ID)
	}

	return &Storage{fs: jsonfs, lastID: lastID}, nil
}

func MustNew(config jsonfile.Config) *Storage {
	storage, err := New(config)
	if err != nil {
		panic(err)
	}

	return storage
}

func (s *Storage) SaveTask(_ context.Context, task *domain.Task) (*domain.Task, error) {
	tasks, err := s.fs.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	s.lastID++
	task.ID = s.lastID
	tasks[task.ID] = fromTask(task)

	err = s.fs.Save(tasks)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	return task, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task *domain.Task) error {
	tasks, err := s.loadIfExist(ctx, task.ID)
	if err != nil {
		return err
	}

	tasks[task.ID] = fromTask(task)

	err = s.fs.Save(tasks)
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, tid uint64) error {
	tasks, err := s.loadIfExist(ctx, tid)
	if err != nil {
		return err
	}

	delete(tasks, tid)

	err = s.fs.Save(tasks)
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, tid uint64) (*domain.Task, error) {
	tasks, err := s.loadIfExist(ctx, tid)
	if err != nil {
		return nil, err
	}

	return toTask(tasks[tid]), nil
}

func (s *Storage) ListAll(_ context.Context) ([]*domain.Task, error) {
	tasks, err := s.fs.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	list := make([]*domain.Task, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, toTask(task))
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("%s error: %w", name, domain.ErrEmptyTasks)
	}

	return list, nil
}

func (s *Storage) ListByStatus(_ context.Context, status domain.Status) ([]*domain.Task, error) {
	tasks, err := s.fs.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	list := make([]*domain.Task, 0)

	for _, task := range tasks {
		if task.Status == status {
			list = append(list, toTask(task))
		}
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("%s error: %w", name, domain.ErrEmptyTasks)
	}

	return list, nil
}
