package game

import (
	"fmt"
	"strconv"

	"github.com/Maycon-Santos/go-snake-backend/uuid"
)

type Matches interface {
	Add(playersLimit int) (Match, error)
	GetMatchByID(id string) (Match, error)
	GetMatchByOwnerID(ownerID string) (Match, error)
	DeleteByID(id string)
}

type matches struct {
	matches map[string]Match
}

func NewMatches() Matches {
	return &matches{
		matches: make(map[string]Match),
	}
}

func (r matches) Add(playersLimit int) (Match, error) {
	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	idStr := strconv.FormatUint(*id, 10)

	match := NewMatch(idStr, playersLimit)

	r.matches[idStr] = match

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
		owner := match.GetOwner()

		if owner == nil {
			continue
		}

		if owner.GetID() == ownerID {
			return match, nil
		}
	}

	return nil, fmt.Errorf("matches: There is no match with id %s owner", ownerID)
}

func (r matches) DeleteByID(id string) {
	delete(r.matches, id)
}
