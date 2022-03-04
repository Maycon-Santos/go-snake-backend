package http

import (
	"encoding/json"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/server/auth"
	"github.com/julienschmidt/httprouter"
)

type signUpRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpResponseResult struct {
	AccessToken string `json:"access_token"`
	RefresToken string `json:"refresh_token"`
}

func SignUpHandler(container container.Container) httprouter.Handle {
	var (
		env                process.Env
		cache              cache.Client
		accountsRepository db.AccountsRepository
	)

	container.Retrieve(&env, &cache, &accountsRepository)

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		var requestBody signUpRequestBody

		err := json.NewDecoder(request.Body).Decode(&requestBody)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		usernameExists, err := accountsRepository.CheckUsernameExists(request.Context(), requestBody.Username)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		if usernameExists {
			responseBody := responseSchema{
				Success: false,
				Type:    TYPE_USERNAME_EXISTS,
				Message: "username already in use",
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		passwordHash, err := auth.GeneratePasswordHash(requestBody.Password)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		accountID, err := accountsRepository.Save(request.Context(), requestBody.Username, passwordHash)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		token, err := auth.CreateToken(
			env.JWT.ExpiresIn,
			env.JWT.RefreshExpiresIn,
			env.JWT.Secret,
			env.JWT.RefreshSecret,
			accountID,
		)
		if err != nil {
			handleError(request.Context(), err)
		}

		if err = auth.CreateAuth(request.Context(), cache, accountID, token); err != nil {
			handleError(request.Context(), err)
		}

		responseBody := responseSchema{
			Success: true,
			Result: signUpResponseResult{
				AccessToken: token.AccessToken,
				RefresToken: token.RefreshToken,
			},
		}
		if err := makeResponse(request.Context(), writer, responseBody); err != nil {
			handleError(request.Context(), err)
		}
	}
}
