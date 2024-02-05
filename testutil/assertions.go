package testutil

import (
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/validation"
	"github.com/golang-jwt/jwt/v5"
)

type GenericTypeToResponseFunction func(interface{}) interface{}

func AssertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatalf("expected error but didn't get one")
	}

	if got != want {
		t.Errorf("got error %v, want %v", got, want)
	}
}

func AssertEvent(t testing.TB, got events.InterfaceEvent, want events.InterfaceEvent) {
	t.Helper()

	if reflect.DeepEqual(got, events.InterfaceEvent{}) {
		t.Fatalf("didn't publish event")
	}

	if got.EventID != want.EventID {
		t.Errorf("got Event ID %v want %v", got.EventID, want.EventID)
	}

	if got.AggregateID != want.AggregateID {
		t.Errorf("got Aggregate ID %v want %v", got.AggregateID, want.AggregateID)
	}

	if !reflect.DeepEqual(got.Payload, want.Payload) {
		t.Errorf("got payload %v want %v", got.Payload, want.Payload)
	}
}

func AssertResponseBody[T, V any](t testing.TB, body io.Reader, data T, ConversionFunction GenericTypeToResponseFunction) {
	t.Helper()

	got, err := validation.ValidateBody[V](body)
	if err != nil {
		t.Fatalf("didn't recieve valid response: %v", err)
	}

	var want V
	var ok bool
	if want, ok = ConversionFunction(data).(V); !ok {
		t.Fatalf("couldn't cast output of ConversionFunction to %T", want)
	}

	AssertEqual(t, got, want)
}

func AssertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func AssertNil(t testing.TB, got interface{}) {
	t.Helper()

	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse httperrors.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Error != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Error, expetedError.Error())
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}

func AssertValidResponse(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("invalid response body: %v", err)
	}
}

func AssertJWT(t testing.TB, jwtString string, secretKey []byte, wantId int) {
	t.Helper()

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		t.Fatalf("error verifying JWT: %v", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		t.Fatalf("did not get subject in JWT, expected %v", wantId)
	}

	gotId, err := strconv.Atoi(subject)
	if err != nil {
		t.Fatalf("did not get customer ID for subject, got %v", subject)
	}

	if gotId != wantId {
		t.Errorf("got customer id %v want %v", subject, wantId)
	}
}
