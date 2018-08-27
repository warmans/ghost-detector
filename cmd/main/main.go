package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/warmans/ghost-detector/pkg/input"
	"github.com/warmans/ghost-detector/pkg/input/button"
	"github.com/warmans/ghost-detector/pkg/output/console"
	"github.com/warmans/ghost-detector/pkg/output/creepy"
	"github.com/warmans/ghost-detector/pkg/output/gauge"
	"go.uber.org/zap"
)

var (
	inputRate  = flag.Int("input.rate", 100, "Read from the input every N milliseconds")
	inputName  = flag.String("input.name", "rand", "source of entropy rand, light")
	outputName = flag.String("output.name", "console", "how to output the result of the input")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()

	// possibly fail if the env is not correct.
	verifyEnv()

	logger.Info("Init device")
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()

	logger.Info("Init inputs")
	var inpt input.Reader
	switch *inputName {
	case "light":
		logger.Fatal("not implemented")
	case "linear":
		inpt = input.NewLinearReader(time.Millisecond * time.Duration(*inputRate))
	default:
		inpt = input.NewRandomReader(time.Millisecond * time.Duration(*inputRate))
	}

	button.New(
		getPin(5, rpio.Input, rpio.Low),
		func() {
			logger.Info("Triggered shutdown")
			cmd := exec.Command("poweroff")
			if err := cmd.Run(); err != nil {
				logger.Fatal("Failed to poweroff", zap.Error(err))
			}
		},
	)

	logger.Info("Init outputs")
	switch *outputName {
	case "creepy":
		//register console output
		creepyOut := creepy.New(2, inpt, os.Stdin)
		defer creepyOut.Close()
	case "physical":
		// register gauge output
		guageOut := gauge.New(inpt, getPin(19, rpio.Pwm, rpio.Low))
		defer guageOut.Close()
		//Register LED output
		ledOut := gauge.New(inpt, getPin(13, rpio.Pwm, rpio.Low))
		defer ledOut.Close()
	default:
		//register console output
		consoleOut := console.New(inpt)
		defer consoleOut.Close()
	}

	logger.Info("Ready")

	// await term
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		logger.Info("Shutting down")
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
