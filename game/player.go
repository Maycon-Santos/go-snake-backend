package game

type Player interface {
	setRoom(room *room)
	GetID() string
	GetName() string
	IsReady() bool
	IsAlive() bool
	MoveToRight()
	MoveToLeft()
	MoveToTop()
	MoveToBottom()
}

type tile struct {
	x int
	y int
}

type player struct {
	id      string
	name    string
	isAlive bool
	isReady bool
	tiles   []tile
	room    *room
}

func NewPlayer(id, name string) Player {
	return &player{
		id:   id,
		name: name,
		tiles: []tile{{
			x: 0,
			y: 0,
		}},
	}
}

func (p player) setRoom(room *room) {
	p.room = room
}

func (p player) GetID() string {
	return p.id
}

func (p player) GetName() string {
	return p.name
}

func (p player) IsReady() bool {
	return p.isReady
}

func (p player) IsAlive() bool {
	return p.isAlive
}

func (p player) MoveToRight() {
	p.tiles = append(
		[]tile{{
			p.tiles[0].x + 1,
			p.tiles[0].y,
		}},
		p.tiles[:len(p.tiles)-1]...,
	)

	if p.tiles[0].x > p.room.tiles.horizontal {
		p.tiles[0].x = 0
	}
}

func (p player) MoveToLeft() {
	p.tiles = append(
		[]tile{{
			p.tiles[0].x - 1,
			p.tiles[0].y,
		}},
		p.tiles[:len(p.tiles)-1]...,
	)

	if p.tiles[0].x < 0 {
		p.tiles[0].x = p.room.tiles.horizontal
	}
}

func (p player) MoveToTop() {
	p.tiles = append(
		[]tile{{
			p.tiles[0].x,
			p.tiles[0].y - 1,
		}},
		p.tiles[:len(p.tiles)-1]...,
	)

	if p.tiles[0].y < 0 {
		p.tiles[0].y = p.room.tiles.vertical
	}
}

func (p player) MoveToBottom() {
	p.tiles = append(
		[]tile{{
			p.tiles[0].x,
			p.tiles[0].y + 1,
		}},
		p.tiles[:len(p.tiles)-1]...,
	)

	if p.tiles[0].y > p.room.tiles.vertical {
		p.tiles[0].y = 0
	}
}
