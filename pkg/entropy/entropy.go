package entropy

import (
	"time"
	"math/rand"
)

type Rander interface {
	Intn(max int) int
}

func NewRand() *Rand {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
	return &Rand{}
}

type Rand struct {
}

func (r *Rand) Intn(max int) int {
	return rand.Intn(max)
}
