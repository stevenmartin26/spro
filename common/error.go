package common

type ErrType int

const (
	ErrUnexpectedError ErrType = iota
	ErrInvalidInput
	ErrUnauthorized
	ErrEntityNotFound
	ErrEntityAlreadyExists
	ErrTooManyAttempts
)

type CustomError struct {
	ErrType ErrType
	Message string
	Details []string
}

func NewCustomError(errType ErrType, message string, details ...string) *CustomError {
	return &CustomError{
		ErrType: errType,
		Message: message,
		Details: details,
	}
}

func (c *CustomError) Error() string {
	return c.Message
}
