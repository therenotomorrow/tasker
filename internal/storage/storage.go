package storage

import (
	"context"
	"fmt"

	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/pkg/jsonfile"
)

const name = "Storage"

type Storage struct {
	engine *jsonfile.JSONFile[Tasks]
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

	return &Storage{engine: jsonfs, lastID: lastID}, nil
}

func MustNew(config jsonfile.Config) *Storage {
	storage, err := New(config)
	if err != nil {
		panic(err)
	}

	return storage
}

func (s *Storage) LastID() uint64 {
	return s.lastID
}

func (s *Storage) SaveTask(_ context.Context, task *domain.Task) (*domain.Task, error) {
	tasks, err := s.engine.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	s.lastID++
	task.ID = s.lastID
	tasks[task.ID] = fromTask(task)

	err = s.engine.Save(tasks)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	return task, nil
}

func (s *Storage) UpdateTask(_ context.Context, task *domain.Task) error {
	tasks, err := s.engine.Load()
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	tasks[task.ID] = fromTask(task)

	err = s.engine.Save(tasks)
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	return nil
}

func (s *Storage) DeleteTask(_ context.Context, task *domain.Task) error {
	tasks, err := s.engine.Load()
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	delete(tasks, task.ID)

	err = s.engine.Save(tasks)
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}

	return nil
}

func (s *Storage) GetByID(_ context.Context, tid uint64) (*domain.Task, error) {
	tasks, err := s.engine.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	if _, ok := tasks[tid]; !ok {
		return nil, fmt.Errorf("%s error: %w", name, domain.ErrTaskNotFound)
	}

	return toTask(tasks[tid]), nil
}

func (s *Storage) ListAll(_ context.Context) ([]*domain.Task, error) {
	tasks, err := s.engine.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	list := make([]*domain.Task, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, toTask(task))
	}

	return list, nil
}

func (s *Storage) ListByStatus(_ context.Context, status domain.Status) ([]*domain.Task, error) {
	tasks, err := s.engine.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	list := make([]*domain.Task, 0)

	for _, task := range tasks {
		domTask := toTask(task)

		if domTask.Status == status {
			list = append(list, domTask)
		}
	}

	return list, nil
}
