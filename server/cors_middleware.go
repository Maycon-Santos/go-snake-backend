package server

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/julienschmidt/httprouter"
)

func CORSMiddleware(container container.Container) func(next httprouter.Handle) httprouter.Handle {
	var (
		env process.Env
	)

	err := container.Retrieve(&env)
	if err != nil {
		log.Fatal(err)
	}

	return func(next httprouter.Handle) httprouter.Handle {
		return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			header := writer.Header()
			header.Set("Access-Control-Allow-Origin", env.AccessControlAllowOrigin)
			header.Set("Access-Control-Allow-Headers", env.AccessControlAllowHeaders)
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))

			if request.Method == "OPTIONS" {
				http.Error(writer, "No Content", http.StatusNoContent)
				return
			}

			next(writer, request, params)
		}
	}
}
