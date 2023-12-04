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
	ErrUnsupportedType      = NewValidationError("validator doesn't support this type")
)

type ErrInvalidRequestField struct {
	msg string
}

func (e *ErrInvalidRequestField) Error() string {
	return e.msg
}

func (e *ErrInvalidRequestField) As(err any) bool {
	if validationError, ok := err.(**ValidationError); ok {
		(*validationError).msg = e.msg

		return true
	}
	return false
}

func NewErrInvalidRequestField(msg string) *ErrInvalidRequestField {
	return &ErrInvalidRequestField{"object has invalid fields:\n" + msg}
}

type ErrInvalidArrayElement struct {
	msg string
	err error
}

func (e *ErrInvalidArrayElement) Error() string {
	return e.msg + e.err.Error()
}

func (e *ErrInvalidArrayElement) Unwrap() error {
	return e.err
}

func NewErrInvalidArrayElement(err error) *ErrInvalidArrayElement {
	return &ErrInvalidArrayElement{"array has invalid elements:\n", err}
}
