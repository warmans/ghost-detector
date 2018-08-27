package button

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

func New(pin rpio.Pin, callback func()) *Push {
	pin.Input()
	pin.PullUp()
	b := &Push{pin: pin, callback: callback}
	go b.start()
	return b
}

type Push struct {
	pin      rpio.Pin
	callback func()
}

func (b *Push) start() {
	pushed := false
	for {
		if b.pin.Read() == rpio.Low {
			if pushed == false {
				b.callback()
			}
			pushed = true
		} else {
			pushed = false
		}
		time.Sleep(time.Millisecond * 10)
	}
}
