package routes

import (
	"encoding/json"
	"math/rand"
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
}

func SignUpHandler(container container.Container) httprouter.Handle {
	var (
		env                process.Env
		cache              cache.Client
		accountsRepository db.AccountsRepository
		skinsRepository    db.SkinsRepository
	)

	container.Retrieve(&env, &cache, &accountsRepository, &skinsRepository)

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		var requestBody signUpRequestBody

		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusUnprocessableEntity,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_PAYLOAD_INVALID,
					Message: "payload is invalid",
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		if responseType, err := validateSignUpFields(requestBody); err != nil {
			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusForbidden,
				},
				Body: responseBody{
					Success: false,
					Type:    responseType,
					Message: err.Error(),
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		usernameExists, err := accountsRepository.CheckUsernameExists(request.Context(), requestBody.Username)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		if usernameExists {
			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusUnauthorized,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_ACCOUNT_USERNAME_EXISTS,
					Message: "username already in use",
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
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

		colors, errColors := skinsRepository.GetAllColors(request.Context())
		if errColors != nil {
			handleError(request.Context(), errColors)
		}

		patterns, errPatterns := skinsRepository.GetAllPatterns(request.Context())
		if errPatterns != nil {
			handleError(request.Context(), errPatterns)
			return
		}

		if errColors == nil && errPatterns == nil {
			err = skinsRepository.SetAccountSkin(
				request.Context(),
				accountID,
				colors[rand.Intn(len(colors)-1)].ID,
				patterns[rand.Intn(len(patterns)-1)].ID,
			)
			if err != nil {
				handleError(request.Context(), err)
				return
			}
		}

		token, err := auth.CreateToken(
			env.JWT.ExpiresIn,
			env.JWT.Secret,
			accountID,
		)
		if err != nil {
			handleError(request.Context(), err)
		}

		if err = auth.CreateAuth(request.Context(), cache, accountID, token); err != nil {
			handleError(request.Context(), err)
		}

		response := responseConfig{
			Header: responseHeader{
				Status: http.StatusCreated,
			},
			Body: responseBody{
				Success: true,
				Result: signUpResponseResult{
					AccessToken: token.AccessToken,
				},
			},
		}
		if err := makeResponse(request.Context(), writer, response); err != nil {
			handleError(request.Context(), err)
		}
	}
}

func validateSignUpFields(requestBody signUpRequestBody) (responseType, error) {
	if errType, err := usernameValidator.Validate(requestBody.Username); err != nil {
		return usernameResponseErrors[errType], err
	}

	if errType, err := passwordValidator.Validate(requestBody.Password); err != nil {
		return passwordResponseErrors[errType], err
	}

	return TYPE_UNKNOWN, nil
}
