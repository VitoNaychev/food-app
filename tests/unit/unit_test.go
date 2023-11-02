package unittest

import (
	"os"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
)

var testEnv handlers.TestEnv

func TestMain(m *testing.M) {
	testEnv = handlers.LoadTestEnviornment()

	code := m.Run()
	os.Exit(code)
}
