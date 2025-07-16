package errors

import "github.com/pkg/errors"

type customError struct {
	errorType     ErrorType
	errorMessage  string
	originalError error
}

// defined methods on CustomError
func (err customError) Message() string {
	return err.errorMessage
}

func (err customError) Code() uint32 {
	return errorTypeToStatusCodeMap[err.errorType]
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
