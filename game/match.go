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
	OnStart(fn func())
	Ready()
	MatchState
}

type match struct {
	ID              string
	playersLimit    int
	owner           Player
	players         []Player
	foods           []Food
	ticker          GameTicker
	onStartHandlers []func()
	onStartSync     sync.Mutex
	foodsSync       sync.Mutex
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

func (m *match) Ready() {
	everyoneIsReady := true

	for _, player := range m.GetPlayers() {
		if !player.IsReady() {
			everyoneIsReady = false
			break
		}
	}

	if everyoneIsReady {
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

	for _, fn := range m.onStartHandlers {
		m.onStartSync.Lock()
		fn()
		m.onStartSync.Unlock()
	}

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

	for _, food := range m.foods {
		food.Summon()
		m.ticker.OnTick(food.CheckWasEaten, 1)
	}
}
