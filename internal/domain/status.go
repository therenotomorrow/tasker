package domain

type Status string

const (
	StatusTodo     Status = "todo"
	StatusDone     Status = "done"
	StatusProgress Status = "progress"
)

func NewStatus(raw string) (Status, error) {
	s := Status(raw)
	switch s {
	case StatusTodo, StatusDone, StatusProgress:
		return s, nil
	default:
		return "", ErrInvalidStatus
	}
}

func AllStatus() []Status {
	return []Status{StatusTodo, StatusProgress, StatusDone}
}
