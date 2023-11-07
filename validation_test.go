package validation_test

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/VitoNaychev/validation"
)

type DummyRequest struct {
	S string `validate:"required,max=20,min=10"`
	I int    `validate:"required"`
}

type IncorrectDummyRequest struct {
	S int
	I string
}

type UnkownFieldsDummyRequest struct {
	S string `validate:"required,max=20,min=10"`
	I int    `validate:"required"`
	F float32
	B []byte
}

type PhoneNumberRequest struct {
	PhoneNumber string `validate:"required,phonenumber"`
}

func TestValidateBody(t *testing.T) {
	t.Run("returns ErrNoBody on no body", func(t *testing.T) {
		_, err := validation.ValidateBody[DummyRequest](nil)

		assertError(t, err, validation.ErrNoBody)
	})

	t.Run("returns ErrEmptyBody on empty body", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})

		_, err := validation.ValidateBody[DummyRequest](body)

		assertError(t, err, validation.ErrEmptyBody)
	})

	t.Run("returns ErrEmptyJSON on empty JSON", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{}`))

		_, err := validation.ValidateBody[DummyRequest](body)

		assertError(t, err, validation.ErrEmptyJSON)
	})

	t.Run("returns ErrIncorrectRequestType on incorrect request type", func(t *testing.T) {
		incorrectDummyRequest := IncorrectDummyRequest{
			S: 10,
			I: "Hello, World!",
		}

		body := newRequestBody(incorrectDummyRequest)

		_, err := validation.ValidateBody[DummyRequest](body)

		assertError(t, err, validation.ErrIncorrectRequestType)
	})

	t.Run("returns ErrIncorrectRequestType on unknown fields", func(t *testing.T) {
		unkownFieldsDummyRequest := UnkownFieldsDummyRequest{
			S: "Hello, World!",
			I: 42,
			F: 3.14,
			B: []byte{'M', 'D', 'M', 'A'},
		}

		body := newRequestBody(unkownFieldsDummyRequest)

		_, err := validation.ValidateBody[DummyRequest](body)

		assertError(t, err, validation.ErrIncorrectRequestType)
	})

	t.Run("returns ErrInvalidRequestField on invalid fields", func(t *testing.T) {
		invalidDummyRequest := DummyRequest{
			S: "Hello,",
			I: 10,
		}

		body := newRequestBody(invalidDummyRequest)

		_, err := validation.ValidateBody[DummyRequest](body)

		assertError(t, err, validation.ErrInvalidRequestField)
	})

	t.Run("parses request body on valid request", func(t *testing.T) {
		wantDummyRequest := DummyRequest{
			S: "Hello, World!",
			I: 10,
		}

		body := newRequestBody(wantDummyRequest)

		gotDummyRequest, err := validation.ValidateBody[DummyRequest](body)

		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}

		if !reflect.DeepEqual(gotDummyRequest, wantDummyRequest) {
			t.Errorf("got %v want %v", gotDummyRequest, wantDummyRequest)
		}
	})

	t.Run("returns ErrInvalidRequestField on invalid phone number", func(t *testing.T) {
		invalidPhoneNumberRequest := PhoneNumberRequest{
			PhoneNumber: "+359 88 4444 abc",
		}

		body := newRequestBody(invalidPhoneNumberRequest)

		_, err := validation.ValidateBody[PhoneNumberRequest](body)

		assertError(t, err, validation.ErrInvalidRequestField)
	})

	t.Run("parses phone number on valid request", func(t *testing.T) {
		phoneNumberRequest := PhoneNumberRequest{
			PhoneNumber: "+359 88 4444 321",
		}

		body := newRequestBody(phoneNumberRequest)

		gotPhoneNumberRequest, err := validation.ValidateBody[PhoneNumberRequest](body)

		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}

		if !reflect.DeepEqual(gotPhoneNumberRequest, phoneNumberRequest) {
			t.Errorf("got %v want %v", gotPhoneNumberRequest, phoneNumberRequest)
		}
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %v want %v", got, want)
	}
}

func newRequestBody(data interface{}) io.Reader {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(data)

	return body
}
