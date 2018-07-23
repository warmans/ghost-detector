package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/warmans/ghost-detector/pkg/words"
	"syscall"
	"os/signal"
	"time"
)

var (
	prefixLen = flag.Int("prefix", 2, "prefix length in words")
	randomize = flag.Bool("randomize", false, "use RNG instead of sensor input")
	wordFrequency = flag.Int("word-frequency", 1, "print a word every N seconds")
)

func main() {
	flag.Parse() // Parse command-line flags.

	var entropy words.Entropy
	if *randomize {
		entropy = words.NewRandEntropy()
	} else {
		entropy = words.NewRandEntropy()
	}

	fmt.Println("Reading input...")
	chain := words.NewChain(*prefixLen)        // Initialize a new Chain.
	chain.Build(os.Stdin)                      // Build chains from standard input.

	fmt.Println("Detecting...")
	chainOut := chain.Generate(time.Duration(*wordFrequency), entropy) // Generate text.
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
