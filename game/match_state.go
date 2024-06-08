package game

import "sync"

type matchStatus string

const (
	StatusOnHold  = matchStatus("ON_HOLD")
	StatusRunning = matchStatus("RUNNING")
)

type MatchState interface {
	UpdateState(input MatchStateInput)
	OnUpdateState(fn func())
	GetMap() Map
	GetFoodsLimit() int
	GetStatus() matchStatus
}

type Tiles struct {
	Horizontal int
	Vertical   int
}

type Map struct {
	Tiles Tiles
}

type matchState struct {
	status           matchStatus
	_map             Map
	foodsLimit       int
	onUpdateHandlers []func()
	sync             sync.Mutex
}

type MapInput struct {
	Tiles *Tiles
}

type MatchStateInput struct {
	Status     *matchStatus
	Map        *MapInput
	FoodsLimit *int
}

func NewMatchState() MatchState {
	return &matchState{}
}

func (ms *matchState) UpdateState(input MatchStateInput) {
	if input.Status != nil {
		ms.status = *input.Status
	}

	if input.Map != nil {
		if input.Map.Tiles != nil {
			ms._map.Tiles = *input.Map.Tiles
		}
	}

	if input.FoodsLimit != nil {
		ms.foodsLimit = *input.FoodsLimit
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

func (ms *matchState) GetMap() Map {
	return ms._map
}

func (ms *matchState) GetFoodsLimit() int {
	return ms.foodsLimit
}

func (ms *matchState) GetStatus() matchStatus {
	status := ms.status
	return status
}
