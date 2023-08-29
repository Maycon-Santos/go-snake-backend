package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/julienschmidt/httprouter"
)

type createRoomResponseResult struct {
	RoomID *uint64 `json:"room_id"`
}

func CreateRoom(container container.Container) httprouter.Handle {
	var (
		env process.Env
	)

	err := container.Retrieve(&env)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")
		accountUsername := params.ByName("account_username")

		if room, err := roomsRepository.GetByOwnerID(accountID); err == nil {
			roomsRepository.DeleteByID(room.ID)
		}

		roomID, err := roomsRepository.Add(5, game.NewPlayer(accountID, accountUsername))
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		err = makeResponse(context.Background(), writer, responseConfig{
			Body: responseBody{
				Success: true,
				Result: createRoomResponseResult{
					RoomID: roomID,
				},
			},
		})
		if err != nil {
			handleError(request.Context(), err)
			return
		}
	}
}
