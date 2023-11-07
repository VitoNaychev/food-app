package validation

import "errors"

var (
	ErrNoBody               = errors.New("request body is nil")
	ErrEmptyBody            = errors.New("request body is empty")
	ErrEmptyJSON            = errors.New("request JSON is empty")
	ErrIncorrectRequestType = errors.New("request type is incorrect")
	ErrInvalidRequestField  = errors.New("request contains invalid field(s)")
)
