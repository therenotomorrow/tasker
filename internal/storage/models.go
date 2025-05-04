package storage

import (
	"tasker/internal/domain"
	"time"
)

type Task struct {
	ID          uint64        `json:"id"` // pk
	Description string        `json:"description"`
	Status      domain.Status `json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type Tasks map[uint64]*Task

func toTask(model *Task) *domain.Task {
	return &domain.Task{
		ID:          model.ID,
		Description: model.Description,
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func fromTask(entity *domain.Task) *Task {
	return &Task{
		ID:          entity.ID,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
