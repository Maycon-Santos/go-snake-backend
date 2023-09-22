package game

import (
	"encoding/json"
	"math"
	"sync"

	"github.com/Maycon-Santos/go-snake-backend/utils"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type Player interface {
	AddMovement(mv movement)
	SendMessage(message []byte) error
	SetMatch(room Match)
	SetSocket(socket *websocket.Conn)
	ToIncrease(toIncrease uint)
	Increase()
	DieOnPlayerCollision()
	GetID() string
	GetName() string
	GenerateInitialBody(n int)
	Move()
	TeleportCornerScreen()
	PlayerState
}

type messageListener = func(message WrittenMessage)

type player struct {
	id     string
	name   string
	socket *websocket.Conn

	match Match

	messageListeners []messageListener
	sendMessageSync  sync.Mutex

	toIncrease uint

	movementSync sync.Mutex
	movements    []movement
	moving       movement
	lastTail     BodyFragment

	PlayerState
}

type WrittenMessage struct {
	MoveTo string `json:"moveTo,omitempty"`
	Ready  *bool  `json:"ready,omitempty"`
}

type movement int

const (
	MoveUp movement = iota
	MoveDown
	MoveLeft
	MoveRight
)

var (
	horizontalMovements = []movement{MoveLeft, MoveRight}
	VerticalMovements   = []movement{MoveUp, MoveDown}
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
	p.sendMessageSync.Lock()
	defer p.sendMessageSync.Unlock()

	p.socket = socket
	p.startListening()
}

func (p *player) startListening() {
	go (func() {
		for {
			messageType, reader, err := p.socket.NextReader()
			if err != nil {
				// Enviar erros para um chan
				return
			}

			if messageType == websocket.TextMessage {
				message := WrittenMessage{}

				err = json.NewDecoder(reader).Decode(&message)
				if err != nil {
					// Enviar erros para um chan
					return
				}

				p.readMessages(message)
			}
		}
	})()
}

func (p *player) readMessages(message WrittenMessage) {
	switch message.MoveTo {
	case "right":
		p.AddMovement(MoveRight)
	case "left":
		p.AddMovement(MoveLeft)
	case "up":
		p.AddMovement(MoveUp)
	case "down":
		p.AddMovement(MoveDown)
	}

	if p.match.GetStatus() == StatusOnHold {
		if message.Ready != nil && *message.Ready {
			p.UpdateState(PlayerStateInput{
				IsReady: message.Ready,
			})

			if *message.Ready {
				p.match.Ready()
			} else {
				p.match.Unready()
			}
		}
	}
}

// Enviar erros para um chan

func (p *player) SendMessage(message []byte) error {
	p.sendMessageSync.Lock()
	defer p.sendMessageSync.Unlock()

	writer, err := p.socket.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = writer.Write(message)
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
	p.movementSync.Lock()
	defer p.movementSync.Unlock()

	nextMovement := p.moving
	if len(p.movements) > 0 {
		nextMovement = p.movements[0]
	}

	if slices.Contains(horizontalMovements, nextMovement) && slices.Contains(horizontalMovements, mv) {
		return
	}

	if slices.Contains(VerticalMovements, nextMovement) && slices.Contains(VerticalMovements, mv) {
		return
	}

	p.movements = append(p.movements, mv)
}

func (p *player) ToIncrease(toIncrease uint) {
	p.toIncrease += toIncrease
}

func (p *player) Move() {
	if !p.IsAlive() {
		return
	}

	p.movementSync.Lock()
	if len(p.movements) > 0 {
		p.moving, p.movements = p.movements[0], p.movements[1:]
	}
	p.movementSync.Unlock()

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
	case MoveUp:
		newBodyFragment = BodyFragment{
			body[0].X,
			body[0].Y - 1,
		}
	case MoveDown:
		newBodyFragment = BodyFragment{
			body[0].X,
			body[0].Y + 1,
		}
	}

	p.lastTail = body[len(p.GetBody())-1]

	p.UpdateState(PlayerStateInput{
		Body: append([]BodyFragment{newBodyFragment}, body[:len(p.GetBody())-1]...),
	})
}

func (p *player) TeleportCornerScreen() {
	if !p.IsAlive() {
		return
	}

	body := p.GetBody()
	head := body[0]
	tiles := p.match.GetMap().Tiles

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
	if !p.IsAlive() {
		return
	}

	if p.toIncrease > 0 {
		p.toIncrease -= 1

		p.UpdateState(PlayerStateInput{
			Body: append(p.GetBody(), p.lastTail),
		})
	}
}

func (p *player) DieOnPlayerCollision() {
	if !p.IsAlive() {
		return
	}

	head := p.GetBody()[0]

	for _, player := range p.match.GetPlayers() {
		if !player.IsAlive() {
			continue
		}

		for j, bodyFragment := range player.GetBody() {
			if j == 0 && player.GetID() == p.id {
				continue
			}

			collided := bodyFragment.X == head.X && bodyFragment.Y == head.Y

			if collided {
				p.movements = make([]movement, 0)
				p.UpdateState(PlayerStateInput{
					IsAlive: utils.Ptr(false),
				})
			}
		}
	}
}

func (p *player) GenerateInitialBody(n int) {
	xMultiplier := (n % 3) + 1
	yMultiplier := math.Ceil(float64(n+1) / 3)

	xBodyStart := int(16 * xMultiplier)
	yBodyStart := int(9 * yMultiplier)

	p.UpdateState(PlayerStateInput{
		Body: []BodyFragment{
			{X: xBodyStart, Y: yBodyStart},
			{X: xBodyStart - 1, Y: yBodyStart},
			{X: xBodyStart - 2, Y: yBodyStart},
		},
	})
}
