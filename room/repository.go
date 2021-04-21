package room

import "github.com/Maycon-Santos/go-snake-backend/uuid"

type room struct {
	playersLimit uint
}

type repository struct {
	rooms map[*uint64]room
}

func NewRepository() repository {
	return repository{}
}

// Add a new room and return the ID or error
func (r repository) add(room room) (id *uint64, err error) {
	id, err = uuid.Generate()

	if err == nil {
		r.rooms[id] = room
	}

	return
}

// Delete room by id
func (r repository) delete(id *uint64) {
	delete(r.rooms, id)
}
