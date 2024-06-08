package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/julienschmidt/httprouter"
)

type updateSkinRequestBody struct {
	ColorID   string `json:"color_id"`
	PatternID string `json:"pattern_id"`
}

func UpdateSkin(container container.Container) httprouter.Handle {
	var skinsRepository db.SkinsRepository

	err := container.Retrieve(&skinsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")

		var requestBody updateSkinRequestBody

		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusUnprocessableEntity,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_PAYLOAD_INVALID,
					Message: "playload is invalid",
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		if colorExists, err := skinsRepository.CheckColorExists(request.Context(), requestBody.ColorID); !colorExists {
			if err != nil {
				response := responseConfig{
					Header: responseHeader{
						Status: http.StatusInternalServerError,
					},
					Body: responseBody{
						Success: false,
						Type:    TYPE_UNKNOWN,
					},
				}

				if err := makeResponse(request.Context(), writer, response); err != nil {
					handleError(request.Context(), err)
				}

				return
			}

			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusForbidden,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_COLOR_NOT_AVAILABLE,
					Message: "The requested color does not exist",
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		if colorExists, err := skinsRepository.CheckPatternsExists(request.Context(), requestBody.PatternID); !colorExists {
			if err != nil {
				response := responseConfig{
					Header: responseHeader{
						Status: http.StatusInternalServerError,
					},
					Body: responseBody{
						Success: false,
						Type:    TYPE_UNKNOWN,
					},
				}

				if err := makeResponse(request.Context(), writer, response); err != nil {
					handleError(request.Context(), err)
				}

				return
			}

			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusForbidden,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_PATTERN_NOT_AVAILABLE,
					Message: "The requested pattern does not exist",
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		err := skinsRepository.SetAccountSkin(request.Context(), accountID, requestBody.ColorID, requestBody.PatternID)
		if err != nil {
			response := responseConfig{
				Header: responseHeader{
					Status: http.StatusUnprocessableEntity,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_UNKNOWN,
				},
			}

			if err := makeResponse(request.Context(), writer, response); err != nil {
				handleError(request.Context(), err)
			}

			return
		}

		response := responseConfig{
			Header: responseHeader{
				Status: http.StatusOK,
			},
			Body: responseBody{
				Success: true,
			},
		}

		if err := makeResponse(request.Context(), writer, response); err != nil {
			handleError(request.Context(), err)
		}
	}
}
