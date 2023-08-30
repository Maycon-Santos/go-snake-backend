package server

import (
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/server/auth"
	"github.com/Maycon-Santos/go-snake-backend/server/routes"
	"github.com/julienschmidt/httprouter"
)

func newRoutes(container container.Container) *httprouter.Router {
	router := httprouter.New()

	authGetDataMiddleware := auth.GetDataMiddleware(container)

	router.POST("/v1/signin", CORS(routes.SignInHandler(container)))
	router.POST("/v1/signup", CORS(routes.SignUpHandler(container)))
	router.GET("/v1/rooms/create", CORS(authGetDataMiddleware(routes.CreateRoom(container))))
	router.GET("/v1/rooms/connect/:room_id", CORS(authGetDataMiddleware(routes.ConnectRoom(container))))

	return router
}
