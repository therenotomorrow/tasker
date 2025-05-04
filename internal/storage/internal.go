package storage

import (
	"context"
	"fmt"
	"tasker/internal/domain"
)

func (s *Storage) loadIfExist(_ context.Context, tid uint64) (Tasks, error) {
	tasks, err := s.fs.Load()
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", name, err)
	}

	if _, ok := tasks[tid]; !ok {
		return nil, fmt.Errorf("%s error: %w", name, domain.ErrTaskNotFound)
	}

	return tasks, nil
}
