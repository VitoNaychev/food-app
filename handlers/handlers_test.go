package handlers_test

import (
	"os"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/config"
)

var testEnv config.Enviornment

func TestMain(m *testing.M) {
	testEnv = config.LoadEnviornment("../config/test.env")

	code := m.Run()
	os.Exit(code)
}
