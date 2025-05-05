package testkit

type constError string

func (e constError) Error() string {
	return string(e)
}

const (
	ErrDummy         constError = "dummy"
	ErrUnimplemented constError = "unimplemented"

	FailureTest = "failure"
	SuccessTest = "success"
)
