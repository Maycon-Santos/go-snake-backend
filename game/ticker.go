package game

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type GameTicker interface {
	OnTick(fn func(), layer uint)
	Stop()
	Reset()
}

type gameTicker struct {
	ticker *time.Ticker
	done   chan bool
	ticks  map[uint][]func()
	sync   sync.Mutex
}

func NewTicker() GameTicker {
	gt := gameTicker{}

	gt.start()

	return &gt
}

func (gt *gameTicker) start() {
	gt.ticker = time.NewTicker(time.Second / 18)
	gt.done = make(chan bool)

	go func() {
		for {
			select {
			case <-gt.done:
				fmt.Println("done")
				return
			case <-gt.ticker.C:
				gt.sync.Lock()

				keys := make([]uint, 0)
				for k := range gt.ticks {
					keys = append(keys, k)
				}

				sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

				for _, k := range keys {
					for _, fn := range gt.ticks[k] {
						fn()
					}
				}

				gt.sync.Unlock()
			}
		}
	}()
}

func (gt *gameTicker) OnTick(fn func(), layer uint) {
	gt.sync.Lock()

	if gt.ticks[layer] == nil {
		gt.ticks[layer] = make([]func(), 0)
	}

	gt.ticks[layer] = append(gt.ticks[layer], fn)

	gt.sync.Unlock()
}

func (gt *gameTicker) Stop() {
	gt.ticker.Stop()
	gt.done <- true
}

func (gt *gameTicker) Reset() {
	gt.sync.Lock()
	gt.ticks = make(map[uint][]func())
	gt.sync.Unlock()
}
