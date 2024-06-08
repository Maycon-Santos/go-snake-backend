package game

import (
	"log"
	"sync"
)

type PlayerState interface {
	UpdateState(input PlayerStateInput)
	OnUpdateState(fn func())
	IsReady() bool
	IsAlive() bool
	GetBody() []BodyFragment

	OpenBatch()
	CloseBatch()
}

type BodyFragment struct {
	X int
	Y int
}

type playerState struct {
	isAlive          bool
	isReady          bool
	body             []BodyFragment
	onUpdateHandlers []func()

	sync sync.Mutex

	isBatching   bool
	onCloseBatch func()
}

type PlayerStateInput struct {
	IsAlive *bool
	IsReady *bool
	Body    []BodyFragment
}

func newPlayerState() PlayerState {
	return &playerState{}
}

func (ps *playerState) OpenBatch() {
	if ps.isBatching {
		log.Fatal("There is already an open batch")
		return
	}

	ps.isBatching = true
}

func (ps *playerState) CloseBatch() {
	if !ps.isBatching {
		log.Fatal("The batch needs to be opened first")
		return
	}

	ps.isBatching = false

	if ps.onCloseBatch != nil {
		ps.onCloseBatch()
		ps.onCloseBatch = nil
	}
}

func (ps *playerState) UpdateState(input PlayerStateInput) {
	if input.IsReady != nil {
		ps.isReady = *input.IsReady
	}

	if input.IsAlive != nil {
		ps.isAlive = *input.IsAlive
	}

	if input.Body != nil {
		ps.body = input.Body
	}

	if ps.isBatching {
		ps.onCloseBatch = ps.dispatchUpdateEvent
	} else {
		ps.dispatchUpdateEvent()
	}
}

func (ps *playerState) dispatchUpdateEvent() {
	ps.sync.Lock()
	defer ps.sync.Unlock()

	if ps.onUpdateHandlers != nil {
		for _, fn := range ps.onUpdateHandlers {
			fn()
		}
	}
}

func (ps *playerState) OnUpdateState(fn func()) {
	ps.sync.Lock()
	defer ps.sync.Unlock()

	ps.onUpdateHandlers = append(ps.onUpdateHandlers, fn)
}

func (ps *playerState) GetBody() []BodyFragment {
	return ps.body
}

func (ps *playerState) IsReady() bool {
	return ps.isReady
}

func (ps *playerState) IsAlive() bool {
	return ps.isAlive
}
