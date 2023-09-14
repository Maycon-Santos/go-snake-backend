package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/Maycon-Santos/go-snake-backend/utils"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ConnectRoom(container container.Container) httprouter.Handle {
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

		socket, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			handleError(request.Context(), err)
		}

		if player := match.GetPlayerByID(accountID); player != nil {
			currentPlayer := *player
			currentPlayer.SetSocket(socket)
			return
		}

		currentPlayer := game.NewPlayer(accountID, accountUsername)
		currentPlayer.SetSocket(socket)

		match.Enter(currentPlayer)

		currentPlayer.OnUpdateState(func() {
			msgBytes, err := parsePlayerMessage(currentPlayer)
			if err != nil {
				handleError(request.Context(), err)
				return
			}

			match.SendMessage(msgBytes)
		})

		currentPlayer.ReadMessage(func(message game.WrittenMessage) {
			switch message.MoveTo {
			case "right":
				currentPlayer.AddMovement(game.MoveRight)
			case "left":
				currentPlayer.AddMovement(game.MoveLeft)
			case "top":
				currentPlayer.AddMovement(game.MoveTop)
			case "bottom":
				currentPlayer.AddMovement(game.MoveBottom)
			}
		})

		if matchMessageBytes, err := parseMatchMessage(match); err == nil {
			err = currentPlayer.SendMessage(matchMessageBytes)
			if err != nil {
				handleError(request.Context(), err)
			}
		} else {
			handleError(request.Context(), err)
		}

		// BOOTSTRAP ⬇️

		currentPlayer.UpdateState(game.PlayerStateInput{
			IsAlive: utils.Ptr(true),
			Body:    []game.BodyFragment{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		})
	}
}
