package game

import (
	"math/rand"
)

type Food interface {
	SetMatch(match Match)
	Summon()
	CheckWasEaten()
	FoodState
}

type food struct {
	match Match
	foodState
}

func NewFood() Food {
	return &food{}
}

func (f *food) SetMatch(match Match) {
	f.match = match
}

func (f *food) Summon() {
	tiles := f.match.GetMap().Tiles

	newPosition := foodPosition{
		X: rand.Intn(tiles.Horizontal - 1),
		Y: rand.Intn(tiles.Vertical - 1),
	}

	for _, player := range f.match.GetPlayers() {
		for _, bodyFragment := range player.GetBody() {
			if bodyFragment.X == newPosition.X && bodyFragment.Y == newPosition.Y {
				f.Summon()
				return
			}
		}
	}

	f.UpdateState(&foodStateInput{
		position: &newPosition,
	})
}

func (f *food) CheckWasEaten() {
	for _, player := range f.match.GetPlayers() {
		head := player.GetBody()[0]

		if !player.IsAlive() {
			continue
		}

		if head.X != f.position.X {
			continue
		}

		if head.Y != f.position.Y {
			continue
		}

		player.ToIncrease(1)
		f.Summon()
	}
}
