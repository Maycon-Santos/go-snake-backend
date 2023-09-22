package auth

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func GetDataMiddleware(container container.Container) func(next httprouter.Handle) httprouter.Handle {
	var (
		env                process.Env
		accountsRepository db.AccountsRepository
	)

	err := container.Retrieve(&env, &accountsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(next httprouter.Handle) httprouter.Handle {
		return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			tokenStr := request.Header.Get("Token")

			if tokenStr == "" {
				tokenStr = request.URL.Query().Get("token")
			}

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
			if !token.Valid {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				account, err := accountsRepository.GetByID(request.Context(), claims["account_id"].(string))
				if err != nil {
					panic(err)
				}

				if account != nil {
					params = append(params, httprouter.Param{
						Key:   "account_id",
						Value: account.ID,
					})

					params = append(params, httprouter.Param{
						Key:   "account_username",
						Value: account.UserName,
					})
				}
			}

			next(writer, request, params)
		}
	}
}
