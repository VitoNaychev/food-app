package bt_customer_svc

import (
	"encoding/json"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	InitValidate()
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match(`^\+[\d ]+$`, []byte(fl.Field().String()))
	return matched
}

func InitValidate() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("phonenumber", validatePhoneNumber)
}

func ValidateBody(body io.Reader, request interface{}) error {
	if body == nil {
		return ErrNoBody
	}

	var maxRequestSize int64 = 10000
	content, err := io.ReadAll(io.LimitReader(body, maxRequestSize))
	if string(content) == "" {
		return ErrEmptyBody
	}

	if string(content) == "{}" {
		return ErrEmptyJSON
	}

	err = json.Unmarshal(content, request)
	if err != nil {
		return ErrIncorrectRequestType
	}

	err = validate.Struct(request)
	if err != nil {
		return ErrInvalidRequestField
	}

	return nil
}
