package game

import "sync"

type matchStatus string

const (
	StatusOnHold  = matchStatus("STATUS_ON_HOLD")
	StatusRunning = matchStatus("STATUS_RUNNING")
)

type MatchState interface {
	UpdateState(input MatchStateInput)
	OnUpdateState(fn func())
	GetArena() Arena
}

type Tiles struct {
	Horizontal int
	Vertical   int
}

type Arena struct {
	Tiles Tiles
}

type matchState struct {
	status           matchStatus
	arena            Arena
	onUpdateHandlers []func()
	sync             sync.Mutex
}

type ArenaInput struct {
	Tiles *Tiles
}

type MatchStateInput struct {
	Status *matchStatus
	Arena  *ArenaInput
}

func NewMatchState() MatchState {
	return &matchState{}
}

func (ms *matchState) UpdateState(input MatchStateInput) {
	if input.Status != nil {
		ms.status = *input.Status
	}

	if input.Arena != nil {
		if input.Arena.Tiles != nil {
			ms.arena.Tiles = *input.Arena.Tiles
		}
	}

	ms.dispatchUpdateEvent()
}

func (ms *matchState) dispatchUpdateEvent() {
	ms.sync.Lock()
	defer ms.sync.Unlock()

	if ms.onUpdateHandlers != nil {
		for _, fn := range ms.onUpdateHandlers {
			fn()
		}
	}
}

func (ms *matchState) OnUpdateState(fn func()) {
	ms.sync.Lock()
	defer ms.sync.Unlock()

	ms.onUpdateHandlers = append(ms.onUpdateHandlers, fn)
}

func (ms *matchState) GetArena() Arena {
	return ms.arena
}
