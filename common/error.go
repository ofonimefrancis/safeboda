package common

import "fmt"

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserError struct {
	err string
}

func NewUserError(text string, args ...interface{}) *UserError {
	return &UserError{fmt.Sprintf(text, args...)}
}

func (e *UserError) Error() string {
	return e.err
}

var (
	ErrSomethingWentWrong = NewUserError("something went wrong")
)
