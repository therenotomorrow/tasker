package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyDescription Error = "emptyDescription"
	ErrTaskNotFound     Error = "taskNotFound"
	ErrInvalidStatus    Error = "invalidStatus"
	ErrInvalidTaskID    Error = "invalidTaskID"
	ErrTaskAlreadyDone  Error = "taskAlreadyDone"
	ErrEmptyTasks       Error = "emptyTasks"
)
