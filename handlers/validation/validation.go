package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
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

func ValidateBody[T any](body io.Reader) (T, error) {
	var requestStruct T

	if body == nil {
		return requestStruct, handlers.ErrNoBody
	}

	var maxRequestSize int64 = 10000
	content, err := io.ReadAll(io.LimitReader(body, maxRequestSize))
	if string(content) == "" {
		return requestStruct, handlers.ErrEmptyBody
	}

	if string(content) == "{}" {
		return requestStruct, handlers.ErrEmptyJSON
	}

	err = strictUnmarshal(content, &requestStruct)
	if err != nil {
		return requestStruct, handlers.ErrIncorrectRequestType
	}

	err = validate.Struct(requestStruct)
	if err != nil {
		return requestStruct, handlers.ErrInvalidRequestField
	}

	return requestStruct, nil
}

func strictUnmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
