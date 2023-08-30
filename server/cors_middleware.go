package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CORS(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
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
