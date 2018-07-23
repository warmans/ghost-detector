package words

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"math/rand"
	"time"
)

type Entropy interface {
	Intn(max int) int
}

func NewRandEntropy() *RandEntropy {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
	return &RandEntropy{}
}

type RandEntropy struct {
}

func (r *RandEntropy) Intn(max int) int {
	return rand.Intn(max)
}

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		found := false
		for _, suffix := range c.chain[key] {
			if s == suffix {
				found = true
			}
		}
		if !found {
			c.chain[key] = append(c.chain[key], s)
		}
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(frequency time.Duration, ent Entropy) chan string {
	out := make(chan string)
	go func() {
		ticker := time.NewTicker(time.Second * frequency)
		p := c.getStartPoint(ent)
		for {
			for range ticker.C {
				choices := c.chain[p.String()]
				if len(choices) == 0 {
					break
				}
				next := choices[ent.Intn(len(choices))]
				out<-next
				p.Shift(next)
			}
		}
	}()
	return out
}

func (c *Chain) getStartPoint(ent Entropy) Prefix {
	start := ent.Intn(len(c.chain))
	for k := range c.chain {
		start--
		if start == 0 {
			return Prefix(strings.Split(k, " "))
		}
	}
	return make(Prefix, c.prefixLen)
}
