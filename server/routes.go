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
	router.GET("/v1/get_account", corsMiddleware(authGetDataMiddleware(routes.GetAccount(container))))
	router.POST("/v1/match/create", corsMiddleware(authGetDataMiddleware(routes.CreateMatch(container))))
	router.GET("/v1/match/connect/:match_id", corsMiddleware(authGetDataMiddleware(routes.ConnectMatch(container))))
	router.GET("/v1/available_skins", corsMiddleware(routes.AvailableSkins(container)))
	router.POST("/v1/update_skin", corsMiddleware(authGetDataMiddleware(routes.UpdateSkin(container))))

	return router
}
