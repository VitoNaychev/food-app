package bt_customer_svc

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

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
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %v want %v", got, want)
	}
}
