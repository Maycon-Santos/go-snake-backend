package server

import (
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/server/http"
	"github.com/Maycon-Santos/go-snake-backend/server/ws"
	"github.com/julienschmidt/httprouter"
)

func newRoutes(container container.Container) *httprouter.Router {
	router := httprouter.New()

	router.POST("/signin", http.SignInHandler(container))
	router.POST("/signup", http.SignUpHandler(container))
	router.GET("/room", ws.Room)

	return router
}
