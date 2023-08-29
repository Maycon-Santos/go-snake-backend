package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type message string

const (
	moveToRightMessage  = message("move to right")
	MoveToLeftMessage   = message("move to left")
	moveToTopMessage    = message("move to top")
	moveToBottomMessage = message("move to bottom")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ConnectRoom(container container.Container) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")
		accountUsername := params.ByName("account_username")

		roomID, err := strconv.ParseUint(params.ByName("room_id"), 10, 64)
		if err != nil {
			handleError(request.Context(), err)
		}

		currentRoom, err := roomsRepository.GetByID(roomID)
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

		player := game.NewPlayer(accountID, accountUsername)

		currentRoom.Enter(player)

		socket, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			handleError(request.Context(), err)
		}

		for {
			msgType, msg, err := socket.ReadMessage()
			if err != nil {
				handleError(request.Context(), err)
				return
			}

			fmt.Println("Mensagem recebida: ", string(msg))

			err = socket.WriteMessage(msgType, msg)
			if err != nil {
				handleError(request.Context(), err)
				return
			}
		}
	}
}
