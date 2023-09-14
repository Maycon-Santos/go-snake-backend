package game

import "fmt"

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

	m.ticker.OnTick(player.Move)
	m.ticker.OnTick(player.TeleportCornerScreen)
	m.ticker.OnTick(player.Increase)
	m.ticker.OnTick(player.DieOnPlayerCollision)

	return nil
}

func (m *match) Bootstrap() {

}

func (m *match) start() {

}
