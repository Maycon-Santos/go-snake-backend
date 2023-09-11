package game

import (
	"fmt"
	"strconv"

	"github.com/Maycon-Santos/go-snake-backend/uuid"
)

type Matches interface {
	Add(playersLimit int, owner Player) (Match, error)
	GetMatchByID(id string) (Match, error)
	GetMatchByOwnerID(ownerID string) (Match, error)
	DeleteByID(id string)
	// SendMessage()
}

type matches struct {
	matches map[string]Match
}

func NewMatches() Matches {
	return &matches{
		matches: make(map[string]Match),
	}
}

func (r matches) Add(playersLimit int, owner Player) (Match, error) {
	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	idStr := strconv.FormatUint(*id, 10)

	match := NewMatch(idStr, owner, playersLimit)

	r.matches[idStr] = match

	owner.SetMatch(r.matches[idStr])

	return match, nil
}

func (r matches) GetMatchByID(id string) (Match, error) {
	if match, ok := r.matches[id]; ok {
		return match, nil
	}

	return nil, fmt.Errorf("matches: There is no match with id %s", id)
}

func (r matches) GetMatchByOwnerID(ownerID string) (Match, error) {
	for _, match := range r.matches {
		if match.GetOwner().GetID() == ownerID {
			return match, nil
		}
	}

	return nil, fmt.Errorf("matches: There is no match with id %s owner", ownerID)
}

func (r matches) DeleteByID(id string) {
	delete(r.matches, id)
}
