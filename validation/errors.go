package validation

type ValidationError struct {
	msg string
	err error
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg, nil}
}

func (v *ValidationError) Error() string {
	return v.msg
}

var (
	ErrNoBody               = NewValidationError("request body is nil")
	ErrEmptyBody            = NewValidationError("request body is empty")
	ErrEmptyJSON            = NewValidationError("request JSON is empty")
	ErrIncorrectRequestType = NewValidationError("request type is incorrect")
)

type ErrInvalidRequestField struct {
	msg string
}

func (v *ErrInvalidRequestField) Error() string {
	return v.msg
}

func (v *ErrInvalidRequestField) As(err any) bool {
	if validationError, ok := err.(**ValidationError); ok {
		(*validationError).msg = v.msg

		return true
	}
	return false
}

func NewErrInvalidRequestField(msg string) *ErrInvalidRequestField {
	return &ErrInvalidRequestField{msg}
}
