package routes

import (
	"encoding/json"

	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/game"
)

type tilesMessage struct {
	Horizontal int `json:"horizontal"`
	Vertical   int `json:"vertical"`
}

type mapMessage struct {
	Tiles tilesMessage `json:"tiles"`
}

type matchMessage struct {
	ID     string     `json:"id"`
	Status string     `json:"status"`
	Map    mapMessage `json:"map"`
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
	Alive    bool                  `json:"alive"`
}

type playerSkinMessage struct {
	PlayerId string `json:"playerId"`
	Color    string `json:"color"`
	Pattern  string `json:"pattern"`
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
	MatchData    *matchMessage      `json:"match,omitempty"`
	Player       *playerMessage     `json:"player,omitempty"`
	PlayerSkin   *playerSkinMessage `json:"playerSkin,omitempty"`
	RemovePlayer string             `json:"removePlayer,omitempty"`
	Food         *foodMessage       `json:"food,omitempty"`
}

func parseMatchMessage(match game.Match) ([]byte, error) {
	mapTiles := match.GetMap().Tiles

	msg := message{
		MatchData: &matchMessage{
			ID:     match.GetID(),
			Status: string(match.GetStatus()),
			Map: mapMessage{
				Tiles: tilesMessage{
					Horizontal: mapTiles.Horizontal,
					Vertical:   mapTiles.Vertical,
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
			Alive:    player.IsAlive(),
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

func parsePlayerSkin(player game.Player, skin db.Skin) ([]byte, error) {
	msg := message{
		PlayerSkin: &playerSkinMessage{
			PlayerId: player.GetID(),
			Color:    skin.ColorID,
			Pattern:  skin.PatternID,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

func parseRemovePlayer(player game.Player) ([]byte, error) {
	msg := message{
		RemovePlayer: player.GetID(),
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
