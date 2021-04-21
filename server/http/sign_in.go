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

type signInRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponseResult struct {
	AccessToken string `json:"access_token"`
	RefresToken string `json:"refresh_token"`
}

func SignInHandler(container container.Container) httprouter.Handle {
	var (
		env                process.Env
		cache              cache.Client
		accountsRepository db.AccountsRepository
	)

	container.Retrieve(&env, &cache, &accountsRepository)

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		var requestBody signInRequestBody

		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			handleError(request.Context(), err)
			return
		}

		account, err := accountsRepository.Get(request.Context(), requestBody.Username)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		if account == nil {
			responseBody := responseConfig{
				Success: false,
				Type:    TYPE_ACCOUNT_NOT_FOUND,
				Message: "account not found",
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		if err = auth.CompareHashAndPassword(account.Password, requestBody.Password); err != nil {
			responseBody := responseConfig{
				Success: false,
				Type:    TYPE_ACCOUNT_PASSWORD_WRONG,
				Message: "wrong password",
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		token, err := auth.CreateToken(env, account.ID)
		if err != nil {
			handleError(request.Context(), err)
		}

		if err = auth.CreateAuth(request.Context(), cache, account.ID, token); err != nil {
			handleError(request.Context(), err)
		}

		responseBody := responseConfig{
			Success: false,
			Result: signInResponseResult{
				AccessToken: token.AccessToken,
				RefresToken: token.RefreshToken,
			},
		}

		if err = makeResponse(request.Context(), writer, responseBody); err != nil {
			handleError(request.Context(), err)
		}
	}
}
