package handlers

import (
	"net/http"
)

type RouterServer struct {
	http.Handler
}

func NewRouterServer(customerServer http.Handler, addressServer http.Handler) *RouterServer {
	routerServer := new(RouterServer)

	router := http.NewServeMux()
	router.Handle("/customer/", customerServer)
	router.Handle("/customer/address/", addressServer)

	routerServer.Handler = router

	return routerServer
}
