package input

import (
	"time"
	"math/rand"
	"github.com/google/uuid"
	"sync"
)

type Reader interface {
	Register(client chan uint) func()
	Close()
}

func NewRandom(frequency time.Duration) *Random {
	r := &Random{
		close: make(chan struct{}, 0),
	}
	rand.Seed(time.Now().UnixNano())
	go r.start(frequency)
	return r
}

// Random is a fake input implementation.
type Random struct {
	clients sync.Map
	close   chan struct{}
}

func (r *Random) Register(client chan uint) func() {
	id := uuid.New().String()
	r.clients.Store(id, client)
	return func() {
		r.clients.Delete(id)
	}
}

func (r *Random) Close() {
	r.close <- struct{}{}
}

func (r *Random) start(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-ticker.C:
			r.clients.Range(func(key, value interface{}) bool {
				value.(chan uint) <- uint(rand.Intn(100))
				return true
			})
		case <-r.close:
			return //exit
		}
	}
}
