package server

import (
	"fmt"
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
			fmt.Println(env.AllowOrigin)
			writer.Header().Add("Access-Control-Allow-Origin", env.AllowOrigin)
			writer.Header().Add("Access-Control-Allow-Credentials", "true")
			writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

			if request.Method == "OPTIONS" {
				http.Error(writer, "No Content", http.StatusNoContent)
				return
			}

			next(writer, request, params)
		}
	}
}
