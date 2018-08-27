package input

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type Reader interface {
	Register(client chan uint) func()
	Close()
}

// AbstractReader implements the register and close methods of a reader but doesn't actually emit any
// values. Some other struct needs to wrap it and publish
type AbstractReader struct {
	clients sync.Map
	close   chan struct{}
}

func (r *AbstractReader) Register(client chan uint) func() {
	id := uuid.Must(uuid.NewV4()).String()
	r.clients.Store(id, client)
	return func() {
		r.clients.Delete(id)
	}
}

func (r *AbstractReader) Close() {
	r.close <- struct{}{}
}

func NewRandomReader(frequency time.Duration) *RandomReader {
	r := &RandomReader{
		AbstractReader: &AbstractReader{
			close: make(chan struct{}, 0),
		},
	}
	rand.Seed(time.Now().UnixNano())
	go r.start(frequency)
	return r
}

// RandomReader is a fake input implementation.
type RandomReader struct {
	*AbstractReader
}

func (r *RandomReader) start(frequency time.Duration) {
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

func NewLinearReader(frequency time.Duration) *LinearReader {
	r := &LinearReader{
		AbstractReader: &AbstractReader{
			close: make(chan struct{}, 0),
		},
	}
	go r.start(frequency)
	return r
}

// LinearReader is a fake input that just cycles from 0 - 100.
type LinearReader struct {
	*AbstractReader
}

func (r *LinearReader) start(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	var count uint
	for {
		select {
		case <-ticker.C:
			r.clients.Range(func(key, value interface{}) bool {
				value.(chan uint) <- count
				return true
			})
			count++
			if count == 100 {
				count = 0
			}
		case <-r.close:
			return //exit
		}
	}
}
