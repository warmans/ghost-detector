package guage

import (
	"github.com/stianeikeland/go-rpio"
)

const cycleLen = 100

func NewPercentageGuage(pin rpio.Pin) *Percentage {
	p := &Percentage{pin: pin}
	p.pin.Freq(64000)
	p.pin.DutyCycle(0, cycleLen)
	return p
}

type Percentage struct {
	pin rpio.Pin
	val int
}

func (g *Percentage) Write(val int) {
	g.pin.DutyCycle(uint32(val), cycleLen)
}
