package handlers

import (
	"net/http"
)

type RouterServer struct {
	http.Handler
}

func NewRouterServer(restaurantServer, addressServer, hoursServer, menuServer http.Handler) *RouterServer {
	routerServer := new(RouterServer)

	router := http.NewServeMux()
	router.Handle("/restaurant/", restaurantServer)
	router.Handle("/restaurant/address/", addressServer)
	router.Handle("/restaurant/hours/", hoursServer)
	router.Handle("/restaurant/menu/", menuServer)

	routerServer.Handler = router

	return routerServer
}
