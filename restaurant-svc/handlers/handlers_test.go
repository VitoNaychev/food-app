package handlers_test

import (
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/loadenv"
	"github.com/VitoNaychev/food-app/testutil"
)

var testEnv loadenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"SECRET", "EXPIRES_AT", "DBUSER", "DBPASS", "DBNAME"}

	var err error
	testEnv, err = loadenv.LoadEnviornment("../config/test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
