package domain

type Error string

const (
	ErrEmptyDescription Error = "emptyDescription"
	ErrTaskNotFound     Error = "taskNotFound"
	ErrInvalidStatus    Error = "invalidStatus"
	ErrInvalidTaskID    Error = "invalidTaskID"
	ErrTaskAlreadyDone  Error = "taskAlreadyDone"
	ErrEmptyTasks       Error = "emptyTasks"
)

func (e Error) Error() string {
	return string(e)
}
