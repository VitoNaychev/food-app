package testdata

import (
	"time"

	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
)

var loc, _ = time.LoadLocation("Local")
var ZeroedTime = time.Date(2000, time.January, 1, 0, 0, 0, 0, loc)

var ReadyByStr = "23:59"
var ReadyByTime, _ = handlers.ParseTimeAndSetDate(ReadyByStr)
