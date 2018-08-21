package main

import (
	"flag"
	"os"
	"fmt"
	"syscall"
	"os/signal"
	"github.com/warmans/ghost-detector/pkg/input"
	"github.com/stianeikeland/go-rpio"
	"github.com/warmans/ghost-detector/pkg/output/console"
	"time"
	"github.com/warmans/ghost-detector/pkg/output/gauge"
)

var (
	sensorName = flag.String("sensor.name", "rand", "source of entropy rand, light")
)

func main() {
	flag.Parse()

	// possibly fail if the env is not correct.
	verifyEnv()

	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()

	var inpt input.Reader
	switch *sensorName {
	case "light":
		panic("not implemented")
	default:
		inpt = input.NewRandom(time.Second)
	}

	//register console output
	consoleOut := console.NewConsoleOutput(inpt)
	defer consoleOut.Close()

	// register gauge output
	guageOut := gauge.NewPercentageOutput(inpt, getPin(19, rpio.Pwm, rpio.Low))
	defer guageOut.Close()

	fmt.Println("Detecting...")

	// await term
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		fmt.Fprint(os.Stderr, "\n\nShutting down")
		//stop all input
		inpt.Close()
		// stop all output
		return
	}
}

func getPin(num uint8, mode rpio.Mode, state rpio.State) rpio.Pin {
	pin := rpio.Pin(num)
	pin.Mode(mode)
	if state == rpio.Low {
		pin.Low()
	} else {
		pin.High()
	}
	return pin
}

func verifyEnv() {
	if "" == os.Getenv("SUDO_USER") {
		panic("pwm requires sudo to work")
	}
}