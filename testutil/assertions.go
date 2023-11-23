package testutil

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/handlers"
	"github.com/VitoNaychev/errorresponse"
)

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertCreateOrderResponse(t testing.TB, got, want handlers.OrderResponse) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertGetOrderResponse(t testing.TB, got, want []handlers.OrderResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse errorresponse.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %v want %v", got, want)
	}
}
