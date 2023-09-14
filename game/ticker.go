package game

import (
	"sync"
	"time"
)

type GameTicker interface {
	OnTick(fn func(), layer uint)
	Stop()
	Reset()
}

type gameTicker struct {
	ticker   *time.Ticker
	done     chan bool
	ticks    map[uint][]func()
	dataSync sync.Mutex
}

func NewTicker() GameTicker {
	gt := gameTicker{}

	gt.start()

	return &gt
}

func (gt *gameTicker) start() {
	gt.ticker = time.NewTicker(500 * time.Millisecond)
	gt.done = make(chan bool)

	go func() {
		for {
			select {
			case <-gt.done:
				return
			case <-gt.ticker.C:
				gt.dataSync.Lock()

				for _, fns := range gt.ticks {
					for _, fn := range fns {
						fn()
					}
				}

				gt.dataSync.Unlock()
			}
		}
	}()
}

func (gt *gameTicker) OnTick(fn func(), layer uint) {
	gt.dataSync.Lock()
	defer gt.dataSync.Unlock()

	if gt.ticks[layer] == nil {
		gt.ticks[layer] = []func(){}
	}

	gt.ticks[layer] = append(gt.ticks[layer], fn)
}

func (gt *gameTicker) Stop() {
	gt.ticker.Stop()
}

func (gt *gameTicker) Reset() {
	gt.dataSync.Lock()
	gt.ticks = map[uint][]func(){}
	gt.done <- true
	gt.dataSync.Unlock()
	gt.start()
}
