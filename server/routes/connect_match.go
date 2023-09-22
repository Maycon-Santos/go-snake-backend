package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func ConnectMatch(container container.Container) httprouter.Handle {
	var matches game.Matches

	err := container.Retrieve(&matches)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")
		accountUsername := params.ByName("account_username")
		matchID := params.ByName("match_id")

		match, err := matches.GetMatchByID(matchID)
		if err != nil {
			makeResponse(request.Context(), writer, responseConfig{
				Header: responseHeader{
					Status: http.StatusNotFound,
				},
				Body: responseBody{
					Success: false,
					Type:    TYPE_ROOM_NOT_FOUND,
					Message: "room not found",
				},
			})
			return
		}

		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}

		socket, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			handleError(request.Context(), err)
		}

		var currentPlayer game.Player

		if player := match.GetPlayerByID(accountID); player != nil {
			currentPlayer = *player
			currentPlayer.SetSocket(socket)
		} else {
			currentPlayer = game.NewPlayer(accountID, accountUsername)
			currentPlayer.SetSocket(socket)

			match.OnStart(func() {
				for _, food := range match.GetFoods() {
					food.OnUpdateState(func() {
						msgBytes, err := parseFoodMessage(food)
						if err != nil {
							handleError(request.Context(), err)
						}

						err = match.SendMessage(msgBytes)
						if err != nil {
							handleError(request.Context(), err)
						}
					})
				}
			})

			currentPlayer.OnUpdateState(func() {
				msgBytes, err := parsePlayerMessage(currentPlayer)
				if err != nil {
					handleError(request.Context(), err)
					return
				}

				match.SendMessage(msgBytes)
			})

			if err = match.Enter(currentPlayer); err != nil {
				handleError(request.Context(), err)
				return
			}
		}

		socket.SetCloseHandler(func(code int, text string) (err error) {
			if match.GetStatus() == game.StatusOnHold {
				match.RemovePlayer(currentPlayer)

				removePlayerMessageBytes, err := parseRemovePlayer(currentPlayer)
				if err != nil {
					handleError(request.Context(), err)
				}

				if match.SendMessage(removePlayerMessageBytes); err != nil {
					handleError(request.Context(), err)
				}
			}

			if len(match.GetPlayers()) == 0 {
				matches.DeleteByID(matchID)
			}

			return
		})

		matchMessageBytes, err := parseMatchMessage(match)
		if err != nil {
			handleError(request.Context(), err)
		}

		if err = currentPlayer.SendMessage(matchMessageBytes); err != nil {
			handleError(request.Context(), err)
		}

		if match.GetStatus() != game.StatusRunning {
			currentPlayerMessageBytes, err := parsePlayerMessage(currentPlayer)
			if err != nil {
				handleError(request.Context(), err)
			}

			for _, player := range match.GetPlayers() {
				playerMessageBytes, err := parsePlayerMessage(player)
				if err != nil {
					handleError(request.Context(), err)
				}

				if err = currentPlayer.SendMessage(playerMessageBytes); err != nil {
					handleError(request.Context(), err)
				}

				if player.GetID() != currentPlayer.GetID() {
					if err = player.SendMessage(currentPlayerMessageBytes); err != nil {
						handleError(request.Context(), err)
					}
				}
			}
		}

		for _, food := range match.GetFoods() {
			foodMessageBytes, err := parseFoodMessage(food)
			if err != nil {
				handleError(request.Context(), err)
			}

			if err = match.SendMessage(foodMessageBytes); err != nil {
				handleError(request.Context(), err)
			}
		}
	}
}
