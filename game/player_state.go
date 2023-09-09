package game

type PlayerState interface {
	UpdateState(input PlayerStateInput)
	OnUpdateState(fn func())
	IsReady() bool
	IsAlive() bool
	GetBody() []BodyFragment
}

type BodyFragment struct {
	X int
	Y int
}

type playerState struct {
	isAlive         bool
	isReady         bool
	Body            []BodyFragment
	onUpdateHandler func()
}

type PlayerStateInput struct {
	IsAlive *bool
	IsReady *bool
	Body    []BodyFragment
}

func newPlayerState() PlayerState {
	return &playerState{}
}

func (ps *playerState) UpdateState(input PlayerStateInput) {
	if input.IsReady != nil {
		ps.isReady = *input.IsReady
	}

	if input.IsAlive != nil {
		ps.isAlive = *input.IsAlive
	}

	if input.Body != nil {
		ps.Body = input.Body
	}

	ps.dispatchUpdateEvent()
}

func (ps *playerState) dispatchUpdateEvent() {
	if ps.onUpdateHandler != nil {
		ps.onUpdateHandler()
	}
}

func (ps *playerState) OnUpdateState(fn func()) {
	ps.onUpdateHandler = fn
}

func (ps *playerState) GetBody() []BodyFragment {
	return ps.Body
}

func (ps *playerState) IsReady() bool {
	return ps.isReady
}

func (ps *playerState) IsAlive() bool {
	return ps.isAlive
}
