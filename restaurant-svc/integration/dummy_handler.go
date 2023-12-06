package integration

import "net/http"

func DummyHandlerFunc(w http.ResponseWriter, r *http.Request) {}

var DummyHandler = http.HandlerFunc(DummyHandlerFunc)
