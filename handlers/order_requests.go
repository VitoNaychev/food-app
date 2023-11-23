package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func NewCancelOrderRequest(jwt string, cancelOrderRequest CancelOrderRequest) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(cancelOrderRequest)

	request, _ := http.NewRequest(http.MethodPost, "/order/cancel/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewCreateOrderRequest(jwt string, createOrderRequestBody CreateOrderRequest) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createOrderRequestBody)

	request, _ := http.NewRequest(http.MethodPost, "/order/new/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewGetCurrentOrdersRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/order/current/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func NewGetAllOrdersRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
	request.Header.Add("Token", jwt)

	return request
}
