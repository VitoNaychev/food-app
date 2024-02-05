package models

import "time"

var loc, _ = time.LoadLocation("Local")
var ZeroTime = time.Date(2000, time.January, 1, 0, 0, 0, 0, loc)
