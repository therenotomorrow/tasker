package usecases

import (
	"strconv"
	"strings"

	"github.com/therenotomorrow/tasker/internal/domain"
)

func (use *UseCases) validateDescription(description string) (string, error) {
	description = strings.TrimSpace(description)

	if description == "" {
		return "", domain.ErrEmptyDescription
	}

	return description, nil
}

func (use *UseCases) validateTaskID(tid string) (uint64, error) {
	taskID, err := strconv.ParseUint(tid, 10, 64)
	if err != nil {
		return 0, domain.ErrInvalidTaskID
	}

	return taskID, nil
}

func (use *UseCases) validateStatus(status string) (domain.Status, error) {
	stat, err := domain.NewStatus(status)
	if err != nil {
		return "", domain.ErrInvalidStatus
	}

	return stat, nil
}
