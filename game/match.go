package game

import (
	"fmt"

	"github.com/Maycon-Santos/go-snake-backend/utils"
)

type Lobby struct {
	IsReady *bool `json:"is_ready,omitempty"`
}

type Message struct {
	Lobby Lobby `json:"lobby,omitempty"`
}

type Match interface {
	SendMessage(message []byte) error
	GetID() string
	GetOwner() Player
	GetPlayers() []Player
	GetPlayerByID(id string) *Player
	Enter(player Player) error
	Bootstrap()
	MatchState
}

type match struct {
	ID           string
	playersLimit int
	owner        Player
	players      []Player
	ticker       GameTicker
	MatchState
}

func NewMatch(id string, playersLimit int) Match {
	return &match{
		ID:           id,
		playersLimit: playersLimit,
		players:      []Player{},
		ticker:       NewTicker(),
		MatchState:   NewMatchState(),
	}
}

func (m *match) SendMessage(message []byte) (err error) {
	for _, player := range m.GetPlayers() {
		err = player.SendMessage(message)
	}

	return
}

func (m *match) playersLen() int {
	return len(m.players) + 1
}

func (m *match) GetID() string {
	return m.ID
}

func (m *match) GetOwner() Player {
	return m.owner
}

func (m *match) GetPlayers() []Player {
	if m.owner == nil {
		return nil
	}

	return append(m.players, m.owner)
}

func (m *match) GetPlayerByID(id string) *Player {
	for _, player := range m.GetPlayers() {
		if player.GetID() == id {
			return &player
		}
	}
	return nil
}

func (m *match) Enter(player Player) error {
	player.SetMatch(m)

	if m.owner == nil {
		m.owner = player
	} else if m.playersLen() < int(m.playersLimit) {
		m.players = append(m.players, player)
	} else {
		return fmt.Errorf("The match already has the maximum number of players (%d)", m.playersLen())
	}

	player.ReadMessage(func(message WrittenMessage) {
		switch message.MoveTo {
		case "right":
			player.AddMovement(MoveRight)
		case "left":
			player.AddMovement(MoveLeft)
		case "top":
			player.AddMovement(MoveTop)
		case "bottom":
			player.AddMovement(MoveBottom)
		}

		if message.Ready != nil && *message.Ready {
			player.UpdateState(PlayerStateInput{
				IsReady: utils.Ptr(true),
			})

			everyoneIsReady := true

			for _, player := range m.GetPlayers() {
				if !player.IsReady() {
					everyoneIsReady = false
					break
				}
			}

			if everyoneIsReady {
				m.start()
			}
		}
	})

	return nil
}

func (m *match) Bootstrap() {

}

func (m *match) start() {
	m.ticker.Reset()

	// Criar sistema de camadas
	for _, player := range m.GetPlayers() {
		player.UpdateState(PlayerStateInput{
			IsAlive: utils.Ptr(true),
			Body:    []BodyFragment{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		})

		m.ticker.OnTick(player.Move, 0)
		m.ticker.OnTick(player.TeleportCornerScreen, 0)
		m.ticker.OnTick(player.Increase, 2)
		m.ticker.OnTick(player.DieOnPlayerCollision, 2)
	}
}
