package testkit

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const (
	ErrDummy         ConstError = "dummy"
	ErrUnimplemented ConstError = "unimplemented"

	FailureTest = "failure"
	SuccessTest = "success"
)
