package integration

import (
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/testutil"
)

var env appenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"SECRET", "DBUSER", "DBPASS", "DBNAME"}

	var err error
	env, err = appenv.LoadEnviornment("../test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
