package bt_customer_svc

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

type DummyRequest struct {
	S string `validate:"required,max=20,min=10"`
	I int    `validate:"required"`
}

type IncorrectDummyRequest struct {
	S int
	I string
}

type PhoneNumberRequest struct {
	PhoneNumber string `validate:"required,phonenumber"`
}

func TestValidateBody(t *testing.T) {
	t.Run("returns Bad Request on no body", func(t *testing.T) {
		var dummyRequest DummyRequest
		err := ValidateBody(nil, &dummyRequest)

		assertError(t, err, ErrNoBody)
	})

	t.Run("returns Bad Request on empty body", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})

		var dummyRequest DummyRequest
		err := ValidateBody(body, &dummyRequest)

		assertError(t, err, ErrEmptyBody)
	})

	t.Run("returns Bad Request on empty JSON", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{}`))

		var dummyRequest DummyRequest
		err := ValidateBody(body, &dummyRequest)

		assertError(t, err, ErrEmptyJSON)
	})

	t.Run("returns Bad Request on incorrect request type", func(t *testing.T) {
		incorrectDummyRequest := IncorrectDummyRequest{
			S: 10,
			I: "Hello, World!",
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(incorrectDummyRequest)

		var dummyRequest DummyRequest
		err := ValidateBody(body, &dummyRequest)

		assertError(t, err, ErrIncorrectRequestType)
	})

	t.Run("returns Bad Request on invalid fields", func(t *testing.T) {
		invalidDummyRequest := DummyRequest{
			S: "Hello,",
			I: 10,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(invalidDummyRequest)

		var dummyRequest DummyRequest
		err := ValidateBody(body, &dummyRequest)

		assertError(t, err, ErrInvalidRequestField)
	})

	t.Run("returns Accepted on valid request", func(t *testing.T) {
		wantDummyRequest := DummyRequest{
			S: "Hello, World!",
			I: 10,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(wantDummyRequest)

		var gotDummyRequest DummyRequest
		err := ValidateBody(body, &gotDummyRequest)

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

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(invalidPhoneNumberRequest)

		var gotDummyRequest DummyRequest
		err := ValidateBody(body, &gotDummyRequest)

		assertError(t, err, ErrInvalidRequestField)
	})

	t.Run("succeeds on valid phone number", func(t *testing.T) {
		phoneNumberRequest := PhoneNumberRequest{
			PhoneNumber: "+359 88 4444 321",
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(phoneNumberRequest)

		var gotPhoneNumberRequest PhoneNumberRequest
		err := ValidateBody(body, &gotPhoneNumberRequest)

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
