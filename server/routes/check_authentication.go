package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func CheckAuthentication(container container.Container) httprouter.Handle {
	var (
		env                process.Env
		cache              cache.Client
		accountsRepository db.AccountsRepository
	)

	err := container.Retrieve(&env, &cache, &accountsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		tokenStr := request.Header.Get("Token")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.JWT.Secret), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if token.Valid {
			writer.WriteHeader(http.StatusOK)
			return
		}

		writer.WriteHeader(http.StatusUnauthorized)
	}
}
