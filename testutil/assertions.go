package testutil

import (
	"reflect"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/handlers"
	"github.com/VitoNaychev/errorresponse"
)

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

func AssertErrorResponse(t testing.TB, got, want errorresponse.ErrorResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %v want %v", got, want)
	}
}
