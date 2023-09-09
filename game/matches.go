package game

import (
	"fmt"

	"github.com/Maycon-Santos/go-snake-backend/uuid"
)

type Matches interface {
	Add(playersLimit int, owner Player) (Match, error)
	GetMatchByID(id uint64) (Match, error)
	GetMatchByOwnerID(ownerID string) (Match, error)
	DeleteByID(id uint64)
	// SendMessage()
}

type matches struct {
	matches map[uint64]Match
}

func NewMatches() Matches {
	return &matches{
		matches: make(map[uint64]Match),
	}
}

func (r matches) Add(playersLimit int, owner Player) (Match, error) {
	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	match := NewMatch(*id, owner, playersLimit)

	r.matches[*id] = match

	owner.SetMatch(r.matches[*id])

	return match, nil
}

func (r matches) GetMatchByID(id uint64) (Match, error) {
	if match, ok := r.matches[id]; ok {
		return match, nil
	}

	return nil, fmt.Errorf("matches: There is no match with id %d", id)
}

func (r matches) GetMatchByOwnerID(ownerID string) (Match, error) {
	for _, match := range r.matches {
		if match.GetOwner().GetID() == ownerID {
			return match, nil
		}
	}

	return nil, fmt.Errorf("matches: There is no match with id %s owner", ownerID)
}

func (r matches) DeleteByID(id uint64) {
	delete(r.matches, id)
}
