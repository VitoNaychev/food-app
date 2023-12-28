package integration

import (
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/loadenv"
	"github.com/VitoNaychev/food-app/testutil"
)

var env loadenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"DBUSER", "DBPASS", "DBNAME"}

	var err error
	env, err = loadenv.LoadEnviornment("../test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
