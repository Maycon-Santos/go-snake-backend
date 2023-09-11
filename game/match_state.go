package game

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
			ms.arena.Tiles = *input.Arena.Tiles
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

func (ms *matchState) GetArena() Arena {
	return ms.arena
}
