package game

import (
	"fmt"
	"sync"

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
	GetFoods() []Food
	Enter(player Player) error
	RemovePlayer(player Player)
	OnStart(fn func())
	Ready()
	Unready()
	MatchState
}

type match struct {
	ID           string
	playersLimit int

	owner        Player
	players      []Player
	foods        []Food
	playersReady int

	ticker          GameTicker
	onStartHandlers []func()

	onStartSync sync.Mutex
	foodsSync   sync.Mutex

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
		if err != nil {
			// Enviar erros para um chan
			continue
		}
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

func (m *match) GetFoods() []Food {
	m.foodsSync.Lock()
	defer m.foodsSync.Unlock()

	return m.foods
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

	return nil
}

func (m *match) RemovePlayer(player Player) {
	if m.owner == player {
		m.owner = nil

		if len(m.players) >= 1 {
			m.owner = m.players[0]
			m.players = m.players[1:]
		}

		return
	}

	for i, p := range m.players {
		if player == p {
			m.players = append(m.players[:i], m.players[i+1:]...)
		}
	}
}

func (m *match) Unready() {
	m.playersReady -= 1

	if m.playersReady < 0 {
		m.playersReady = 0
	}
}

func (m *match) Ready() {
	m.playersReady += 1

	if m.playersReady == len(m.GetPlayers()) {
		m.UpdateState(MatchStateInput{
			Status: utils.Ptr(StatusRunning),
		})

		m.start()
	}
}

func (m *match) OnStart(fn func()) {
	m.onStartSync.Lock()
	defer m.onStartSync.Unlock()

	m.onStartHandlers = append(m.onStartHandlers, fn)
}

func (m *match) start() {
	m.ticker.Reset()

	m.foodsSync.Lock()
	m.foods = make([]Food, 0, m.GetFoodsLimit())
	m.foodsSync.Unlock()

	for i := 0; i < m.GetFoodsLimit(); i++ {
		food := NewFood()
		food.SetMatch(m)
		m.foods = append(m.foods, food)
	}

	m.onStartSync.Lock()
	for _, fn := range m.onStartHandlers {
		fn()
	}
	m.onStartSync.Unlock()

	for i, player := range m.GetPlayers() {
		player := player

		player.OnDie(func() {
			for _, p := range m.GetPlayers() {
				if p.IsAlive() {
					return
				}
			}

			m.end()
		})

		player.UpdateState(PlayerStateInput{
			IsReady: utils.Ptr(false),
			IsAlive: utils.Ptr(true),
		})

		player.GenerateInitialBody(i)

		m.playersReady = 0

		m.ticker.OnTick(func() {
			player.OpenBatch()
			player.Move()
			player.TeleportCornerScreen()
		}, 0)

		m.ticker.OnTick(func() {
			player.Increase()
			player.DieOnPlayerCollision()
			player.CloseBatch()
		}, 2)
	}

	for _, food := range m.foods {
		food.Summon()
		m.ticker.OnTick(food.CheckWasEaten, 1)
	}
}

func (m *match) end() {
	m.UpdateState(MatchStateInput{
		Status: utils.Ptr(StatusOnHold),
	})

	m.foods = make([]Food, 0)

	for _, player := range m.GetPlayers() {
		player.Reset()
	}
}
