package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// type message string

// const (
// 	moveToRightMessage  = message("move to right")
// 	MoveToLeftMessage   = message("move to left")
// 	moveToTopMessage    = message("move to top")
// 	moveToBottomMessage = message("move to bottom")
// )

type bodyFragment struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type tiles struct {
	Horizontal int `json:"horizontal"`
	Vertical   int `json:"vertical"`
}

type Player struct {
	ID       string         `json:"id"`
	Username string         `json:"username"`
	Body     []bodyFragment `json:"body"`
}

type Arena struct {
	Tiles tiles `json:"tiles"`
}

type Match struct {
	ID    uint64 `json:"id"`
	Arena Arena  `json:"arena"`
}

type message struct {
	Player    *Player `json:"player,omitempty"`
	MatchData *Match  `json:"match,omitempty"`
}

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

		matchID, err := strconv.ParseUint(params.ByName("room_id"), 10, 64)
		if err != nil {
			handleError(request.Context(), err)
		}

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

		var currentPlayer game.Player

		if player := match.GetPlayerByID(accountID); player != nil {
			currentPlayer = *player
		} else {
			currentPlayer = game.NewPlayer(accountID, accountUsername)
		}

		currentPlayer.SetSocket(socket)

		match.Enter(currentPlayer)

		match.OnUpdateState(func() {
			msg := message{
				MatchData: &Match{
					ID: match.GetID(),
					Arena: Arena{
						Tiles: tiles{
							Horizontal: match.GetTiles().Horizontal,
							Vertical:   match.GetTiles().Vertical,
						},
					},
				},
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				handleError(request.Context(), err)
			}

			err = match.SendMessage(msgBytes)
			if err != nil {
				handleError(request.Context(), err)
			}
		})

		currentPlayer.OnUpdateState(func() {
			msg := message{
				Player: &Player{
					ID:       currentPlayer.GetID(),
					Username: currentPlayer.GetName(),
				},
			}

			for _, fragment := range currentPlayer.GetBody() {
				msg.Player.Body = append(msg.Player.Body, bodyFragment{
					X: fragment.X,
					Y: fragment.Y,
				})
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				handleError(request.Context(), err)
			}

			err = match.SendMessage(msgBytes)
			if err != nil {
				handleError(request.Context(), err)
			}
		})

		// BOOTSTRAP ⬇️

		match.UpdateState(game.MatchStateInput{
			Arena: &game.ArenaInput{
				Tiles: &game.Tiles{
					Horizontal: 60,
					Vertical:   60,
				},
			},
		})

		currentPlayer.UpdateState(game.PlayerStateInput{
			Body: []game.BodyFragment{{X: 0, Y: 0}},
		})

		ticker := game.NewTicker()

		for _, player := range match.GetPlayers() {
			ticker.OnTick(player.Move)
		}
	}
}
