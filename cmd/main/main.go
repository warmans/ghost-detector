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
	"github.com/warmans/ghost-detector/pkg/output/creepy"
)

var (
	inputRate = flag.Int("input.rate", 100, "Read from the input every N milliseconds")
	inputName  = flag.String("input.name", "rand", "source of entropy rand, light")
	outputName = flag.String("output.name", "console", "how to output the result of the input")
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
	switch *inputName {
	case "light":
		panic("not implemented")
	case "linear":
		inpt = input.NewLinearReader(time.Millisecond * time.Duration(*inputRate))
	default:
		inpt = input.NewRandomReader(time.Millisecond * time.Duration(*inputRate))
	}

	switch *outputName {
	case "creepy":
		//register console output
		creepyOut := creepy.New(2, inpt, os.Stdin)
		defer creepyOut.Close()
	case "physical":
		// register gauge output
		guageOut := gauge.New(inpt, getPin(19, rpio.Pwm, rpio.Low))
		defer guageOut.Close()
	default:
		//register console output
		consoleOut := console.New(inpt)
		defer consoleOut.Close()
	}

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