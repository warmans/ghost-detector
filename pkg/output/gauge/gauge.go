package gauge

import (
	"github.com/stianeikeland/go-rpio"
	"github.com/warmans/ghost-detector/pkg/input"
)

const cycleLen = 100

func New(input input.Reader, pin rpio.Pin) *Percentage {
	out := &Percentage{
		in:   make(chan uint, 100),
		stop: make(chan struct{}),
		pin:  pin,
	}

	//configure pwm pin
	out.pin.Freq(64000)
	out.pin.DutyCycle(0, cycleLen)

	out.deregister = input.Register(out.in)
	out.start()
	return out
}

type Percentage struct {
	in         chan uint
	deregister func()
	stop       chan struct{}
	pin        rpio.Pin
}

func (o *Percentage) Close() {
	o.deregister()
	close(o.in)
	o.stop <- struct{}{}
}

func (o *Percentage) start() {
	go func() {
		for {
			select {
			case v := <-o.in:
				o.pin.DutyCycle(uint32(v), cycleLen)
			case <-o.stop:
				return
			}
		}
	}()
}
