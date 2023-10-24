package bt_customer_svc

import (
	"encoding/json"
	"io"

	"github.com/asaskevich/govalidator"
)

func ValidateRequest(bodyReader io.Reader, request interface{}) error {
	var maxRequestSize int64 = 10000
	body, err := io.ReadAll(io.LimitReader(bodyReader, maxRequestSize))
	if string(body) == "" {
		return ErrEmptyBody
	}

	if string(body) == "{}" {
		return ErrEmptyJSON
	}

	err = json.Unmarshal(body, request)
	if err != nil {
		return ErrIncorrectRequestType
	}

	valid, _ := govalidator.ValidateStruct(request)
	if !valid {
		return ErrInvalidRequestField
	}

	return nil
}
