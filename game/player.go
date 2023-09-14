package game

import (
	"encoding/json"
	"sync"

	"github.com/Maycon-Santos/go-snake-backend/utils"
	"github.com/gorilla/websocket"
)

type Player interface {
	AddMovement(mv movement)
	ReadMessage(fn messageListener)
	SendMessage(message []byte) error
	SetMatch(room Match)
	SetSocket(socket *websocket.Conn)
	ToIncrease(toIncrease uint)
	Increase()
	DieOnPlayerCollision()
	GetID() string
	GetName() string
	Move()
	TeleportCornerScreen()
	PlayerState
}

type messageListener = func(message WrittenMessage)

type player struct {
	id               string
	name             string
	toIncrease       uint
	lastTail         BodyFragment
	movements        []movement
	moving           movement
	socket           *websocket.Conn
	messageListeners []messageListener
	messageSync      sync.Mutex
	match            Match
	PlayerState
}

type WrittenMessage struct {
	MoveTo string `json:"moveTo"`
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
	p.startListening()
}

func (p *player) startListening() {
	go (func() {
		for {
			messageType, messageBytes, err := p.socket.ReadMessage()
			if err != nil {
				// Enviar erros para um chan
				return
			}

			if messageType == websocket.TextMessage {
				message := WrittenMessage{}

				err = json.Unmarshal(messageBytes, &message)
				if err != nil {
					// Enviar erros para um chan
					return
				}

				for _, listener := range p.messageListeners {
					p.messageSync.Lock()
					listener(message)
					p.messageSync.Unlock()
				}
			}
		}
	})()
}

// Enviar erros para um chan

func (p *player) ReadMessage(fn messageListener) {
	p.messageListeners = append(p.messageListeners, fn)
}

func (p *player) SendMessage(message []byte) error {
	p.messageSync.Lock()
	defer p.messageSync.Unlock()

	writer, err := p.socket.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err = writer.Write(message); err != nil {
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

func (p *player) ToIncrease(toIncrease uint) {
	p.toIncrease += toIncrease
}

func (p *player) Move() {
	p.messageSync.Lock()

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
			body[0].X,
			body[0].Y + 1,
		}
	}

	p.lastTail = body[len(p.GetBody())-1]

	p.messageSync.Unlock()

	p.UpdateState(PlayerStateInput{
		Body: append([]BodyFragment{newBodyFragment}, body[:len(p.GetBody())-1]...),
	})
}

func (p *player) TeleportCornerScreen() {
	body := p.GetBody()
	head := body[0]
	tiles := p.match.GetArena().Tiles

	if head.X >= tiles.Horizontal {
		p.UpdateState(PlayerStateInput{
			Body: append([]BodyFragment{{X: 0, Y: head.Y}}, body[1:]...),
		})
	}

	if head.X < 0 {
		p.UpdateState(PlayerStateInput{
			Body: append([]BodyFragment{{X: tiles.Horizontal - 1, Y: head.Y}}, body[1:]...),
		})
	}

	if head.Y >= tiles.Vertical {
		p.UpdateState(PlayerStateInput{
			Body: append([]BodyFragment{{X: head.X, Y: 0}}, body[1:]...),
		})
	}

	if head.Y < 0 {
		p.UpdateState(PlayerStateInput{
			Body: append([]BodyFragment{{X: head.X, Y: tiles.Vertical - 1}}, body[1:]...),
		})
	}
}

func (p *player) Increase() {
	if p.toIncrease > 0 {
		p.toIncrease -= 1

		p.UpdateState(PlayerStateInput{
			Body: append(p.GetBody(), p.lastTail),
		})
	}
}

func (p *player) DieOnPlayerCollision() {
	head := p.GetBody()[0]

	for _, player := range p.match.GetPlayers() {
		for j, bodyFragment := range player.GetBody() {
			if j == 0 && player.GetID() == p.id {
				continue
			}

			collided := bodyFragment.X == head.X && bodyFragment.Y == head.Y

			if player.IsAlive() && collided {
				p.UpdateState(PlayerStateInput{
					IsAlive: utils.Ptr(false),
				})
			}
		}
	}
}
