package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/server/auth"
	"github.com/Maycon-Santos/go-snake-backend/validator"
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

	err := container.Retrieve(&env, &cache, &accountsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		var requestBody signInRequestBody

		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			responseBody := responseConfig{
				Status:  http.StatusUnprocessableEntity,
				Success: false,
				Type:    TYPE_PAYLOAD_INVALID,
				Message: "playload invalid",
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		if responseType, err := validateFields(requestBody); err != nil {
			responseBody := responseConfig{
				Status:  http.StatusForbidden,
				Success: false,
				Type:    responseType,
				Message: err.Error(),
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

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
				Status:  http.StatusForbidden,
				Success: false,
				Type:    TYPE_ACCOUNT_PASSWORD_WRONG,
				Message: "wrong password",
			}

			if err := makeResponse(request.Context(), writer, responseBody); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		token, err := auth.CreateToken(
			env.JWT.ExpiresIn,
			env.JWT.RefreshExpiresIn,
			env.JWT.Secret,
			env.JWT.RefreshSecret,
			account.ID,
		)
		if err != nil {
			handleError(request.Context(), err)
		}

		if err = auth.CreateAuth(request.Context(), cache, account.ID, token); err != nil {
			handleError(request.Context(), err)
		}

		responseBody := responseConfig{
			Success: true,
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

func validateFields(requestBody signInRequestBody) (responseType, error) {
	usernameValidator := validator.
		Field("username").
		Required().
		MinLen(4).
		MaxLen(15).
		NoContains([]string{" "})

		// TODO: Não precisa de min, nem max length
	usernameResponseErrors := map[string]responseType{
		validator.Required: TYPE_MISSING_USERNAME,
		validator.MinLen:   TYPE_BELOW_LEN_USERNAME,
		validator.MaxLen:   TYPE_ABOVE_LEN_USERNAME,
	}

	if errType, err := usernameValidator.Validate(requestBody.Username); err != nil {
		return usernameResponseErrors[errType], err
	}

	passwordValidator := validator.
		Field("username").
		Required().
		MinLen(6).
		MaxLen(25)

		// TODO: Não precisa de min, nem max length
	passwordResponseErrors := map[string]responseType{
		validator.Required: TYPE_MISSING_PASSWORD,
		validator.MinLen:   TYPE_BELOW_LEN_PASSWORD,
		validator.MaxLen:   TYPE_ABOVE_LEN_PASSWORD,
	}

	if errType, err := passwordValidator.Validate(requestBody.Password); err != nil {
		return passwordResponseErrors[errType], err
	}

	return TYPE_UNKNOWN, nil
}
