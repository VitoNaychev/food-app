package handlers

import (
	"net/http"
)

type RouterServer struct {
	http.Handler
}

func NewRouterServer(restaurantServer, addressServer http.Handler) *RouterServer {
	routerServer := new(RouterServer)

	router := http.NewServeMux()
	router.Handle("/restaurant/", restaurantServer)
	router.Handle("/restaurant/address/", addressServer)

	routerServer.Handler = router

	return routerServer
}
