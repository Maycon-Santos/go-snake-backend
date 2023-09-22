package game

import (
	"fmt"
	"strconv"
	"sync"

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
	sync    sync.Mutex
}

func NewMatches() Matches {
	return &matches{
		matches: make(map[string]Match),
	}
}

func (m *matches) Add(playersLimit int) (Match, error) {
	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	idStr := strconv.FormatUint(*id, 10)

	match := NewMatch(idStr, playersLimit)

	m.matches[idStr] = match

	return match, nil
}

func (m *matches) GetMatchByID(id string) (Match, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	if match, ok := m.matches[id]; ok {
		return match, nil
	}

	return nil, fmt.Errorf("matches: There is no match with id %s", id)
}

func (m *matches) GetMatchByOwnerID(ownerID string) (Match, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	for _, match := range m.matches {
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

func (m *matches) DeleteByID(id string) {
	m.sync.Lock()
	defer m.sync.Unlock()

	delete(m.matches, id)
}
