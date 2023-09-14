package routes

import (
	"encoding/json"

	"github.com/Maycon-Santos/go-snake-backend/game"
)

type tilesMessage struct {
	Horizontal int `json:"horizontal"`
	Vertical   int `json:"vertical"`
}

type arenaMessage struct {
	Tiles tilesMessage `json:"tiles"`
}

type matchMessage struct {
	ID    string       `json:"id"`
	Arena arenaMessage `json:"arena"`
}

type bodyFragmentMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type playerMessage struct {
	ID       string                `json:"id"`
	Username string                `json:"username"`
	Body     []bodyFragmentMessage `json:"body"`
	Ready    bool                  `json:"ready"`
}

type foodPositionMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type foodMessage struct {
	ID       string              `json:"id"`
	Position foodPositionMessage `json:"position"`
}

type message struct {
	MatchData *matchMessage  `json:"match,omitempty"`
	Player    *playerMessage `json:"player,omitempty"`
	Food      *foodMessage   `json:"food,omitempty"`
}

func parseMatchMessage(match game.Match) ([]byte, error) {
	arenaTiles := match.GetArena().Tiles

	msg := message{
		MatchData: &matchMessage{
			ID: match.GetID(),
			Arena: arenaMessage{
				Tiles: tilesMessage{
					Horizontal: arenaTiles.Horizontal,
					Vertical:   arenaTiles.Vertical,
				},
			},
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

func parsePlayerMessage(player game.Player) ([]byte, error) {
	msg := message{
		Player: &playerMessage{
			ID:       player.GetID(),
			Username: player.GetName(),
			Ready:    player.IsReady(),
			Body:     make([]bodyFragmentMessage, 0),
		},
	}

	for _, fragment := range player.GetBody() {
		msg.Player.Body = append(msg.Player.Body, bodyFragmentMessage{
			X: fragment.X,
			Y: fragment.Y,
		})
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

func parseFoodMessage(food game.Food) ([]byte, error) {
	foodPosition := food.GetPosition()

	msg := message{
		Food: &foodMessage{
			Position: foodPositionMessage{
				X: foodPosition.X,
				Y: foodPosition.Y,
			},
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}
