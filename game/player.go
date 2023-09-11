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
	DieOnPlayerCollision()
	GetID() string
	GetName() string
	Move()
	PlayerState
}

type messageListener = func(message WrittenMessage)

type player struct {
	id               string
	name             string
	movements        []movement
	moving           movement
	socket           *websocket.Conn
	messageListeners []messageListener
	readMessageSync  sync.Mutex
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
					p.readMessageSync.Lock()
					listener(message)
					p.readMessageSync.Unlock()
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
	p.readMessageSync.Lock()

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

	tiles := p.match.GetArena().Tiles

	// Passar isso para um método chamado TeleportCornerScreen

	if newBodyFragment.X >= tiles.Horizontal {
		newBodyFragment.X = 0
	}

	if newBodyFragment.X < 0 {
		newBodyFragment.X = tiles.Horizontal - 1
	}

	if newBodyFragment.Y >= tiles.Vertical {
		newBodyFragment.Y = 0
	}

	if newBodyFragment.Y < 0 {
		newBodyFragment.Y = tiles.Vertical - 1
	}

	p.UpdateState(PlayerStateInput{
		Body: append([]BodyFragment{newBodyFragment}, body[:len(p.GetBody())-1]...),
	})
	p.readMessageSync.Unlock()
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
