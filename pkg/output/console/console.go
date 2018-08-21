package console

import (
	"github.com/warmans/ghost-detector/pkg/input"
	"fmt"
)

func New(input input.Reader) *Output {
	out := &Output{
		in: make(chan uint, 100),
		stop: make(chan struct{}),
	}
	out.deregister = input.Register(out.in)
	out.start()
	return out
}

type Output struct {
	in         chan uint
	deregister func()
	stop       chan struct{}
}

func (o *Output) Close() {
	o.deregister()
	close(o.in)
	o.stop <- struct{}{}
}

func (o *Output) start() {
	go func() {
		for {
			select {
			case v := <-o.in:
				fmt.Printf("Observed value: %d\n", v)
			case <-o.stop:
				return
			}
		}
	}()
}