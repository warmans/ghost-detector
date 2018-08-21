package creepy

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"github.com/warmans/ghost-detector/pkg/input"
	"github.com/warmans/ghost-detector/pkg/util"
)

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

// New returns a new Chain with prefixes of prefixLen words.
func New(prefixLen int, entropy input.Reader, chainData io.Reader) *Chain {

	ch := &Chain{
		chain:     make(map[string][]string),
		prefixLen: prefixLen,

		in:   make(chan uint, 100),
		stop: make(chan struct{}),
	}

	// build the data chains
	ch.build(chainData)

	ch.deregister = entropy.Register(ch.in)
	ch.start()

	return ch
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int

	in         chan uint
	deregister func()
	stop       chan struct{}
}

func (c *Chain) Close() {
	c.deregister()
	close(c.in)
	c.stop <- struct{}{}
}

// build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) build(r io.Reader) {
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

func (c *Chain) start() {
	go func() {
		var p Prefix
		for {
			select {
			case v := <-c.in:

				// we need at least one value to select a start point
				if p == nil {
					p = c.getStartPoint(int(v))
					continue
				}

				// then just keep outputting forever.
				choices := c.chain[p.String()]
				if len(choices) == 0 {
					break
				}
				next := choices[util.Normalize(int(v), 100, len(choices))]
				p.Shift(next)

				// output something
				fmt.Printf("%s ", next)

			case <-c.stop:
				return
			}
		}
	}()
}

func (c *Chain) getStartPoint(startPoint int) Prefix {
	start := util.Normalize(startPoint, 100, len(c.chain))
	for k := range c.chain {
		start--
		if start == 0 {
			return Prefix(strings.Split(k, " "))
		}
	}
	return make(Prefix, c.prefixLen)
}
