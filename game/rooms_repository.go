package game

import (
	"fmt"

	"github.com/Maycon-Santos/go-snake-backend/uuid"
)

type RoomsRepository interface {
	Add(playersLimit int, owner Player) (id *uint64, err error)
	GetByID(id uint64) (*room, error)
	GetByOwnerID(ownerID string) (*room, error)
	DeleteByID(id uint64)
}

type repository struct {
	rooms map[uint64]*room
}

func NewRoomsRepository() RoomsRepository {
	return repository{
		rooms: make(map[uint64]*room),
	}
}

// Add a new room and return the ID or error
func (r repository) Add(playersLimit int, owner Player) (id *uint64, err error) {
	id, err = uuid.Generate()
	if err != nil {
		return nil, err
	}

	r.rooms[*id] = &room{
		ID:           *id,
		playersLimit: playersLimit,
		status:       StatusOnHold,
		owner:        owner,
		players:      []Player{},
	}

	owner.setRoom(r.rooms[*id])

	return id, nil
}

func (r repository) GetByID(id uint64) (*room, error) {
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}

	return nil, fmt.Errorf("roomsRepository: There is no room with id %d", id)
}

func (r repository) GetByOwnerID(ownerID string) (*room, error) {
	for _, room := range r.rooms {
		if room.owner.GetID() == ownerID {
			return room, nil
		}
	}

	return nil, fmt.Errorf("roomsRepository: There is no room with id %s owner", ownerID)
}

// DeleteByID room by id
func (r repository) DeleteByID(id uint64) {
	delete(r.rooms, id)
}
