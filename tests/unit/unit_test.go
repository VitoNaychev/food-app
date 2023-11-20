package unittest

import (
	"os"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/tests"
)

var testEnv tests.TestEnv

func TestMain(m *testing.M) {
	testEnv = tests.LoadTestEnviornment()

	code := m.Run()
	os.Exit(code)
}
