package server

import (
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/server/http"
	"github.com/Maycon-Santos/go-snake-backend/server/ws"
	"github.com/julienschmidt/httprouter"
)

func newRoutes(container container.Container) *httprouter.Router {
	router := httprouter.New()

	router.POST("/v1/signin", http.SignInHandler(container))
	router.POST("/v1/signup", http.SignUpHandler(container))
	router.GET("/v1/room", ws.Room)

	return router
}
