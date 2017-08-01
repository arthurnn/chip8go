package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	//	"fmt"
	"github.com/arthurnn/chip8"
	"github.com/nsf/termbox-go"
)

func main() {

	emulator := chip8.NewChip8()

	//emulator.Output = chip8.NewTermboxDisplay()
	//defer emulator.Output.Close()

	// load game
	rom, err := ioutil.ReadFile("../roms/TETRIS")
	if err != nil {
		log.Fatal(err)
	}

	emulator.LoadRom(rom)

	clock := time.Tick(time.Second / time.Duration(60))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			ev := termbox.PollEvent()

			if ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD {
				stop <- syscall.SIGINT
			}
		}

	}()

loop:
	for {
		select {
		case <-stop:
			break loop
		case <-clock:
			emulator.Run()
			emulator.Output.Render()
		}

		// save key pressed

		//		switch ev := termbox.PollEvent(); ev.Type {
		//		case termbox.EventKey:
		//			switch ev.Key {
		//			case termbox.KeyEsc:
		//				break loop
		//			}
		//		}

	}
}
