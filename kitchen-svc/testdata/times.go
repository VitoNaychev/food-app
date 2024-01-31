package testdata

import (
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
)

var ReadyByStr = "23:59"
var ReadyByTime, _ = handlers.ParseTimeAndSetDate(ReadyByStr)
