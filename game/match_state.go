package game

type matchStatus string

const (
	StatusOnHold  = matchStatus("STATUS_ON_HOLD")
	StatusRunning = matchStatus("STATUS_RUNNING")
)

type MatchState interface {
	UpdateState(input MatchStateInput)
	OnUpdateState(fn func())
	GetTiles() Tiles
}

type Tiles struct {
	Horizontal int
	Vertical   int
}

type Arena struct {
	tiles Tiles
}

type matchState struct {
	tiles           Tiles
	status          matchStatus
	arena           Arena
	onUpdateHandler func()
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
			ms.arena.tiles = *input.Arena.Tiles
		}
	}
	ms.dispatchUpdateEvent()
}

func (ms *matchState) dispatchUpdateEvent() {
	if ms.onUpdateHandler != nil {
		ms.onUpdateHandler()
	}
}

func (ms *matchState) OnUpdateState(fn func()) {
	ms.onUpdateHandler = fn
}

func (ms *matchState) GetTiles() Tiles {
	return ms.tiles
}
