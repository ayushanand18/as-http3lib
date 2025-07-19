package errors

import (
	"net/http"

	"github.com/pkg/errors"
)

type customError struct {
	errorType     ErrorType
	errorMessage  string
	originalError error
}

// defined methods on CustomError
func (err customError) Message() string {
	return err.errorMessage
}

func (err customError) Code() int {
	code, ok := errorTypeToStatusCodeMap[err.errorType]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

func (err customError) String() string {
	return errorTypeToMessageMap[err.errorType]
}

func (err customError) Error() string {
	return err.originalError.Error()
}

// defined methods on ErrorType
func (errType ErrorType) New(message string) error {
	return customError{errorType: errType, originalError: errors.New(message)}
}

func (errType ErrorType) Wrap(err error, message string) error {
	return customError{errorType: errType, originalError: errors.Wrap(err, message)}
}
