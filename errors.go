package bt_customer_svc

import "errors"

var (
	ErrExistingUser         = errors.New("user with this email already exists")
	ErrMissingCustomer      = errors.New("customer doesn't exists")
	ErrMissingToken         = errors.New("missing token")
	ErrInvalidCredentials   = errors.New("invalid user credentials")
	ErrMissingSubject       = errors.New("token does not contain subject field")
	ErrNonIntegerSubject    = errors.New("token subject field is not an integer")
	ErrNoBody               = errors.New("request body is nil")
	ErrEmptyBody            = errors.New("request body is empty")
	ErrEmptyJSON            = errors.New("request JSON is empty")
	ErrIncorrectRequestType = errors.New("request type is incorrect")
	ErrInvalidRequestField  = errors.New("request contains invalid field(s)")
	ErrMissingAddress       = errors.New("address doesn't exists")
	ErrUnathorizedAction    = errors.New("customer does not have permission to perform this action")
)
