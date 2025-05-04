package domain

import "time"

type Task struct {
	ID          uint64
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t Task) IsDone() bool {
	return t.Status == StatusDone
}
