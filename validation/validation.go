package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidationObject interface {
	interface{} | []interface{}
}

var validate *validator.Validate

func init() {
	InitValidate()
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match(`^\+[\d ]+$`, []byte(fl.Field().String()))
	return matched
}

func validateWorkingHours(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match(`^(2[0-3]|[0-1]\d):[0-5]\d$`, []byte(fl.Field().String()))
	return matched
}

func InitValidate() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("phonenumber", validatePhoneNumber)
	validate.RegisterValidation("workinghours", validateWorkingHours)
}

func ValidateBody[T ValidationObject](body io.Reader) (T, error) {
	var requestObject T

	if body == nil {
		return requestObject, ErrNoBody
	}

	var maxRequestSize int64 = 10000
	content, err := io.ReadAll(io.LimitReader(body, maxRequestSize))
	if string(content) == "" {
		return requestObject, ErrEmptyBody
	}

	if string(content) == "{}" {
		return requestObject, ErrEmptyJSON
	}

	err = strictUnmarshal(content, &requestObject)
	if err != nil {
		return requestObject, ErrIncorrectRequestType
	}

	switch reflect.TypeOf(requestObject).Kind() {
	case reflect.Slice, reflect.Array:
		requestArray := reflect.ValueOf(requestObject)
		for i := 0; i < requestArray.Len(); i++ {
			elem := requestArray.Index(i)
			err := ValidateStruct(elem.Interface())
			if err != nil {
				return requestObject, NewErrInvalidArrayElement(err)
			}
		}
	case reflect.Struct:
		err = ValidateStruct(requestObject)
		if err != nil {
			return requestObject, err
		}
	default:
		return requestObject, ErrUnsupportedType
	}

	return requestObject, nil
}

func ValidateStruct(requestStruct interface{}) error {
	err := validate.Struct(requestStruct)
	if err != nil {
		return NewErrInvalidRequestField(err.Error())
	}

	return nil
}

func strictUnmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
