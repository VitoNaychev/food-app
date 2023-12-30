package integrationtest

import (
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"

	"github.com/VitoNaychev/food-app/testutil"
)

var testEnv appenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"SECRET", "EXPIRES_AT", "DBUSER", "DBPASS", "DBNAME"}

	var err error
	testEnv, err = appenv.LoadEnviornment("../test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
