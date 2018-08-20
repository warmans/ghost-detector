package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/warmans/ghost-detector/pkg/words"
	"syscall"
	"os/signal"
	"github.com/warmans/ghost-detector/pkg/entropy"
	"github.com/warmans/ghost-detector/pkg/entropy/sensor"
	"github.com/stianeikeland/go-rpio"
	"time"
	"github.com/warmans/ghost-detector/pkg/output/guage"
)

var (
	prefixLen     = flag.Int("word.prefix", 2, "prefix length in words")
	wordFrequency = flag.Int("word.frequency", 1, "print a word every N seconds")
	sensorName = flag.String("sensor.name", "rand", "source of entropy rand, light")
)

func main() {
	flag.Parse()

	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()

	var ent entropy.Rander
	switch *sensorName {
	case "light":
		ent = sensor.NewLightSensor(getPin(4, rpio.Input, rpio.Low))
	default:
		ent = entropy.NewRand()
	}

	fmt.Println("Reading input...")
	chain := words.NewChain(*prefixLen) // Initialize a new Chain.
	chain.Build(os.Stdin)               // Build chains from standard input.


	//guage
	pin := rpio.Pin(19)
	pin.Mode(rpio.Pwm)
	g := guage.NewPercentageGuage(pin)

	fmt.Println("Detecting...")
	chainOut := chain.Generate(time.Duration(*wordFrequency), ent) // Generate text.
	for {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-c:
			fmt.Fprint(os.Stderr, "\n\nShutting down")
			return
		case out := <-chainOut:
			fmt.Printf("%s ", out)
			g.Write(ent.Intn(100))
		}
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
