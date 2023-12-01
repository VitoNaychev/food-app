package reqbuilder

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func NewRequestWithBody[T any](method string, url string, object T) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(object)

	request, _ := http.NewRequest(method, url, body)
	return request
}

func SetRequestJWT(request *http.Request, jwt string) {
	request.Header.Set("Token", jwt)
}
