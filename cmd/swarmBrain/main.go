package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	brain "github.com/dshills/swarm"
)

func main() {
	swarm := ""
	flag.StringVar(&swarm, "swarm", swarm, "Load a swarm brain")
	flag.Parse()

	if swarm == "" {
		fmt.Println("swarm is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	thebrain, err := brain.Load(swarm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Type /bye to exit")

	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	// Check if the scanner is valid
	if scanner == nil {
		fmt.Println("Failed to create scanner")
		return
	}

	fmt.Print("Question >  ")
	// Read a single line from stdin
	for scanner.Scan() {
		// Get the line as a string
		line := scanner.Text()
		if line == "/bye" {
			os.Exit(0)
		}
		ch := thebrain.Think(line)

		out := <-ch

		for _, r := range out.Results {
			fmt.Println(r)
		}
		fmt.Println(out.String())
		fmt.Println()
		fmt.Print("Question >  ")
	}
}
