package handlers

import (
	"net/http"
)

type RouterServer struct {
	http.Handler
}

func NewRouterServer(deliveryServer *DeliveryServer, locationServer *LocationServer) *RouterServer {
	routerServer := new(RouterServer)

	router := http.NewServeMux()
	router.Handle("/delivery/", deliveryServer)
	router.Handle("/delivery/location/", locationServer)

	routerServer.Handler = router

	return routerServer
}
