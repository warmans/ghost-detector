package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/warmans/ghost-detector/pkg/words"
	"syscall"
	"os/signal"
	"time"
	"github.com/warmans/ghost-detector/pkg/entropy"
	"github.com/warmans/ghost-detector/pkg/entropy/sensor"
	"github.com/warmans/go-rpio"
	"log"
)

var (
	prefixLen     = flag.Int("word.prefix", 2, "prefix length in words")
	wordFrequency = flag.Int("word.frequency", 1, "print a word every N seconds")
	sensorName = flag.String("sensor.name", "rand", "source of entropy rand, light")
	sensorSimulate      = flag.Bool("sensor.simulate", false, "Use a simulated GPIO")
)

func main() {
	flag.Parse()

	device := makeDevice()
	if err := device.Open(); err != nil {
		log.Fatal(err)
	}

	var ent entropy.Rander
	if *sensorName == "light" {
		ent = sensor.NewLightSensor(device.Pin(4, rpio.Output, rpio.PullOff))
	} else {
		ent = entropy.NewRand()
	}

	fmt.Println("Reading input...")
	chain := words.NewChain(*prefixLen) // Initialize a new Chain.
	chain.Build(os.Stdin)               // Build chains from standard input.

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
		}
	}
}


func makeDevice() rpio.Device {
	// use
	var device rpio.Device
	if *sensorSimulate {
		device = rpio.NewPi3Simulator(true)
	} else {
		device = rpio.NewPhysicalDevice()
	}
	return device
}
