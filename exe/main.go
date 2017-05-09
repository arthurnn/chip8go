package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/arthurnn/chip8"
)

func main() {

	emulator := chip8.NewChip8()

	// load game
	rom, err := ioutil.ReadFile("../roms/ibm")
	if err != nil {
		log.Fatal(err)
	}

	emulator.LoadRom(rom)

	for {
		emulator.Run()

		emulator.Output.Render()

		time.Sleep((1000 / 60) * time.Millisecond)

		// save key pressed
	}
}
