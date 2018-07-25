package sensor

import (
	"github.com/warmans/go-rpio"
	"time"
)

const maxChargeTime = 50000

func NewLightSensor(pin *rpio.Pin) *LightSensor {
	sensor := &LightSensor{pin: pin}
	go sensor.start()
	return sensor
}

type LightSensor struct {
	pin          *rpio.Pin
	currentValue int
}

func (s *LightSensor) Intn(max int) int {
	return normalize(s.currentValue, maxChargeTime, max)
}

func (s *LightSensor) start() {
	for {
		count := 0
		s.pin.Output()
		s.pin.Low()
		time.Sleep(time.Millisecond)

		s.pin.Input()
		for s.pin.Read() == rpio.Low {
			time.Sleep(time.Microsecond)
			count += 1
			if count == maxChargeTime {
				break
			}
		}
		s.currentValue = count
	}
}

// normalize take an int with a scale of e.g. 0-1000 and converts it to an int
// in the scale of 0-20.
func normalize(val, oldMax, newMax int) int {
	return int(float64(newMax) / float64(oldMax) * float64(val))
}
