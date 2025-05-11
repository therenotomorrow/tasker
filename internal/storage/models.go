package storage

import (
	"time"

	"github.com/therenotomorrow/tasker/internal/domain"
)

type Task struct {
	ID          uint64    `json:"id"` // pk
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Tasks map[uint64]*Task

func toTask(model *Task) *domain.Task {
	return &domain.Task{
		ID:          model.ID,
		Description: model.Description,
		Status:      domain.Status(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func fromTask(entity *domain.Task) *Task {
	return &Task{
		ID:          entity.ID,
		Description: entity.Description,
		Status:      string(entity.Status),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
