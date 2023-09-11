package server

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/server/auth"
	"github.com/Maycon-Santos/go-snake-backend/server/routes"
	"github.com/julienschmidt/httprouter"
)

func newRoutes(container container.Container) *httprouter.Router {
	var (
		env process.Env
	)

	err := container.Retrieve(&env)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	authGetDataMiddleware := auth.GetDataMiddleware(container)
	corsMiddleware := CORSMiddleware(container)

	router.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Access-Control-Request-Method") != "" {
			header := writer.Header()
			header.Set("Access-Control-Allow-Origin", env.AccessControlAllowOrigin)
			header.Set("Access-Control-Allow-Headers", env.AccessControlAllowHeaders)
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
		}

		writer.WriteHeader(http.StatusNoContent)
	})

	router.POST("/v1/signin", corsMiddleware(routes.SignInHandler(container)))
	router.POST("/v1/signup", corsMiddleware(routes.SignUpHandler(container)))
	router.GET("/v1/check_authentication", corsMiddleware(routes.CheckAuthentication(container)))
	router.POST("/v1/rooms/create", corsMiddleware(authGetDataMiddleware(routes.CreateRoom(container))))
	router.GET("/v1/rooms/connect/:match_id", corsMiddleware(authGetDataMiddleware(routes.ConnectRoom(container))))

	return router
}
