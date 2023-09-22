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
	MatchID string `json:"match_id"`
}

func CreateMatch(container container.Container) httprouter.Handle {
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

		if match, err := matches.GetMatchByOwnerID(accountID); err == nil {
			matches.DeleteByID(match.GetID())
		}

		match, err := matches.Add(5)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		match.UpdateState(game.MatchStateInput{
			Status:     utils.Ptr(game.StatusOnHold),
			FoodsLimit: utils.Ptr(1),
			Map: &game.MapInput{
				Tiles: &game.Tiles{
					Horizontal: 64,
					Vertical:   36,
				},
			},
		})

		match.OnUpdateState(func() {
			msgBytes, err := parseMatchMessage(match)
			if err != nil {
				handleError(request.Context(), err)
				return
			}

			err = match.SendMessage(msgBytes)
			if err != nil {
				handleError(request.Context(), err)
			}
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
