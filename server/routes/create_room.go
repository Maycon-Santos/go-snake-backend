package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/utils"
	"github.com/julienschmidt/httprouter"
)

type createRoomResponseResult struct {
	MatchID uint64 `json:"match_id"`
}

func CreateRoom(container container.Container) httprouter.Handle {
	var (
		env     process.Env
		matches game.Matches
	)

	err := container.Retrieve(&env, &matches)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")
		accountUsername := params.ByName("account_username")

		playerOwner := game.NewPlayer(accountID, accountUsername)

		if match, err := matches.GetMatchByOwnerID(accountID); err == nil {
			matches.DeleteByID(match.GetID())
		}

		match, err := matches.Add(5, playerOwner)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		match.UpdateState(game.MatchStateInput{
			Status: utils.Ptr(game.StatusOnHold),
		})

		err = makeResponse(context.Background(), writer, responseConfig{
			Body: responseBody{
				Success: true,
				Result: createRoomResponseResult{
					MatchID: match.GetID(),
				},
			},
		})
		if err != nil {
			handleError(request.Context(), err)
			return
		}
	}
}
