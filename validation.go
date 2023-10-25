package bt_customer_svc

import (
	"encoding/json"
	"io"

	"github.com/asaskevich/govalidator"
)

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

	valid, _ := govalidator.ValidateStruct(request)
	if !valid {
		return ErrInvalidRequestField
	}

	return nil
}
