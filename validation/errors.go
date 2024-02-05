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
	ErrNoBody          = NewValidationError("message body is nil")
	ErrEmptyBody       = NewValidationError("message body is empty")
	ErrEmptyJSON       = NewValidationError("message JSON is empty")
	ErrUnsupportedType = NewValidationError("validator doesn't support this type")
)

type ErrIncorrectRequestType struct {
	msg string
}

func NewErrIncorrectRequestType(err error) *ErrIncorrectRequestType {
	return &ErrIncorrectRequestType{"message type is incorrect: " + err.Error()}
}

func (e *ErrIncorrectRequestType) Error() string {
	return e.msg
}

func (e *ErrIncorrectRequestType) As(err any) bool {
	if validationError, ok := err.(**ValidationError); ok {
		(*validationError).msg = e.msg

		return true
	}
	return false
}

type ErrInvalidRequestField struct {
	msg string
}

func NewErrInvalidRequestField(msg string) *ErrInvalidRequestField {
	return &ErrInvalidRequestField{"object has invalid fields:\n" + msg}
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

type ErrInvalidArrayElement struct {
	msg string
	err error
}

func NewErrInvalidArrayElement(err error) *ErrInvalidArrayElement {
	return &ErrInvalidArrayElement{"array has invalid elements:\n", err}
}

func (e *ErrInvalidArrayElement) Error() string {
	return e.msg + e.err.Error()
}

func (e *ErrInvalidArrayElement) Unwrap() error {
	return e.err
}
