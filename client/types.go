package client

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

const ClientErrorCode models.ErrorCode = "CLIENT_ERROR"

type ClientError struct {
	Message string
	Code    models.ErrorCode
}

func (e *ClientError) Error() string {
	return e.Message
}

func NewClientError(code models.ErrorCode, msg string) error {
	return errors.WithStack(&ClientError{
		Code:    code,
		Message: msg,
	})
}

func NewClientErrorf(code models.ErrorCode, format string, arg ...any) error {
	msg := fmt.Sprintf(format, arg...)
	return errors.WithStack(&ClientError{
		Code:    code,
		Message: msg,
	})
}

func IsClientError(err error) bool {
	var testError *ClientError
	return errors.As(err, &testError)
}
