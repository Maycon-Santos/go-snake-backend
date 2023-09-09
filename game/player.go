package game

import (
	"github.com/gorilla/websocket"
)

type Player interface {
	AddMovement(mv movement)
	SendMessage(message []byte) error
	SetMatch(room Match)
	SetSocket(socket *websocket.Conn)
	GetID() string
	GetName() string
	Move()
	PlayerState
}

type player struct {
	id        string
	name      string
	movements []movement
	moving    movement
	socket    *websocket.Conn
	match     Match
	PlayerState
}

type movement int

const (
	MoveTop movement = iota
	MoveBottom
	MoveLeft
	MoveRight
)

func NewPlayer(id, name string) Player {
	return &player{
		id:          id,
		name:        name,
		moving:      MoveRight,
		PlayerState: newPlayerState(),
	}
}

func (p *player) SetSocket(socket *websocket.Conn) {
	p.socket = socket
}

func (p *player) SendMessage(message []byte) error {
	err := p.socket.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}

func (p *player) SetMatch(match Match) {
	p.match = match
}

func (p *player) GetID() string {
	return p.id
}

func (p *player) GetName() string {
	return p.name
}

func (p *player) AddMovement(mv movement) {
	p.movements = append(p.movements, mv)
}

func (p *player) Move() {
	if len(p.movements) > 0 {
		p.moving, p.movements = p.movements[0], p.movements[1:]
	}

	body := p.GetBody()

	var newBodyFragment BodyFragment

	switch p.moving {
	case MoveRight:
		newBodyFragment = BodyFragment{
			body[0].X + 1,
			body[0].Y,
		}
	case MoveLeft:
		newBodyFragment = BodyFragment{
			body[0].X - 1,
			body[0].Y,
		}
	case MoveTop:
		newBodyFragment = BodyFragment{
			body[0].X,
			body[0].Y - 1,
		}
	case MoveBottom:
		newBodyFragment = BodyFragment{
			body[0].X + 1,
			body[0].Y,
		}
	}

	if newBodyFragment.X >= p.match.GetTiles().Horizontal {
		newBodyFragment.X = 0
	}

	if newBodyFragment.X < 0 {
		newBodyFragment.X = p.match.GetTiles().Horizontal - 1
	}

	if newBodyFragment.Y >= p.match.GetTiles().Vertical {
		newBodyFragment.Y = 0
	}

	if newBodyFragment.Y < 0 {
		newBodyFragment.Y = p.match.GetTiles().Vertical - 1
	}

	p.UpdateState(PlayerStateInput{
		Body: append([]BodyFragment{newBodyFragment}, body[:len(p.GetBody())-1]...),
	})
}
