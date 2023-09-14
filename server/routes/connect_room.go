package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
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
			p := *player
			p.SetSocket(socket)
			return
		}

		player := game.NewPlayer(accountID, accountUsername)
		player.SetSocket(socket)

		match.Enter(player)

		player.OnUpdateState(func() {
			msgBytes, err := parsePlayerMessage(player)
			if err != nil {
				handleError(request.Context(), err)
				return
			}

			match.SendMessage(msgBytes)
		})

		if matchMessageBytes, err := parseMatchMessage(match); err == nil {
			err = player.SendMessage(matchMessageBytes)
			if err != nil {
				handleError(request.Context(), err)
			}
		} else {
			handleError(request.Context(), err)
		}

		// BOOTSTRAP ⬇️

		// 🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫🎫

		// food := game.NewFood()
		// food.SetMatch(match)

		// food.OnUpdateState(func() {
		// 	msgBytes, err := parseFoodMessage(food)
		// 	if err != nil {
		// 		handleError(request.Context(), err)
		// 	}

		// 	err = match.SendMessage(msgBytes)
		// 	if err != nil {
		// 		handleError(request.Context(), err)
		// 	}
		// })

		// food.Summon()

	}
}
