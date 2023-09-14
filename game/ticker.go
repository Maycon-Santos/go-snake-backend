package game

import (
	"sync"
	"time"
)

type GameTicker interface {
	OnTick(func())
	Stop()
}

type gameTicker struct {
	ticker   *time.Ticker
	ticks    []func()
	dataSync sync.Mutex
}

func NewTicker() GameTicker {
	ticker := time.NewTicker(500 * time.Millisecond)

	gt := gameTicker{
		ticker: ticker,
	}

	go func() {
		for range ticker.C {
			gt.sync(func() {
				for _, fn := range gt.ticks {
					fn()
				}
			})
		}
	}()

	return &gt
}

// copiar essa forma para os states

func (gt *gameTicker) OnTick(fn func()) {
	gt.sync(func() {
		gt.ticks = append(gt.ticks, fn)
	})
}

func (gt *gameTicker) sync(fn func()) {
	gt.dataSync.Lock()
	fn()
	gt.dataSync.Unlock()
}

func (gt *gameTicker) Stop() {
	gt.ticker.Stop()
}
