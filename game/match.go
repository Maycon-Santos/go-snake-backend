package game

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
	Enter(player Player)
	MatchState
}

type match struct {
	ID           string
	playersLimit int
	owner        Player
	players      []Player
	MatchState
}

func NewMatch(id string, owner Player, playersLimit int) Match {
	return &match{
		ID:           id,
		playersLimit: playersLimit,
		owner:        owner,
		players:      []Player{},
		MatchState:   NewMatchState(),
	}
}

func (m *match) SendMessage(message []byte) error {
	for _, player := range m.GetPlayers() {
		if err := player.SendMessage(message); err != nil {
			return err
		}
	}

	return nil
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

func (m *match) Enter(player Player) {
	player.SetMatch(m)

	if m.owner.GetID() == player.GetID() {
		m.owner = player
		return
	}

	if m.playersLen() < int(m.playersLimit) {
		m.players = append(m.players, player)
	}
}

func (m *match) start() {
	// m.UpdateState(MatchStateInput{
	// 	Tiles: &Tiles{
	// 		Horizontal: 60,
	// 		Vertical:   60,
	// 	},
	// })

	// // currentPlayer.UpdateState(game.PlayerStateInput{
	// // 	Body: []game.BodyFragment{{X: 0, Y: 0}},
	// // })

	// ticker := NewTicker()

	// for _, player := range m.GetPlayers() {
	// 	ticker.OnTick(player.Move)
	// }
}
