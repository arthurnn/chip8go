package main

import (
	//	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Graphics struct {
	drawFlag bool
	gfx      [64][32]bool
}

func (g *Graphics) Render() {
	if g.drawFlag {
		// draw here
		//fmt.Println("drawing...")
		g.drawFlag = false

		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				if g.gfx[x][y] {
					fmt.Printf("*")
				} else {
					fmt.Printf(" ")
				}
			}
			fmt.Printf("\n")
		}
	}
}
func (g *Graphics) SetPixel(x, y byte, memory []byte) (collision bool) {
	g.drawFlag = true

	//fmt.Printf("set pixel size %x\n", memory)
	var rx, ry int
	for ry = 0; ry < len(memory); ry++ {
		pixel := memory[ry]
		for rx = 0; rx < 8; rx++ {
			p := 0x80 >> uint(rx)
			if (pixel & byte(p)) > 0 {
				if g.gfx[int(x)+rx][int(y)+ry] {
					g.gfx[int(x)+rx][int(y)+ry] = false
				} else {
					g.gfx[int(x)+rx][int(y)+ry] = true
				}

			}
		}
	}

	return
}

type Chip8 struct {
	//	CPU    cpu
	V      [16]byte
	PC     uint16
	I      uint16
	DTimer byte
	STime  byte

	Memory [4096]byte
	//	// keyboard input

	Output Graphics
}

func main() {

	emulator := Chip8{
		PC: 0x200,
	}

	// load game
	rom, err := ioutil.ReadFile("zero")
	if err != nil {
		log.Fatal(err)
	}

	romArea := emulator.Memory[emulator.PC:]
	copy(romArea, rom)

	for {
		Cycle(&emulator)
		if emulator.DTimer > 0 {
			emulator.DTimer -= 1
		}

		emulator.Output.Render()

		time.Sleep((1000 / 60) * time.Millisecond)

		// save key pressed
	}
}

func Cycle(emulator *Chip8) {
	//fmt.Printf("PC IS %x\n", emulator.PC)

	var opcode uint16
	opcode = uint16(emulator.Memory[emulator.PC])<<8 | uint16(emulator.Memory[emulator.PC+1])

	if 0x1000 == (opcode & 0xF000) {
		// 1NNN	Flow	goto NNN;	Jumps to address NNN.
		emulator.PC = opcode & 0x0FFF

		//		foo := uint16(emulator.Memory[emulator.PC])<<8 | uint16(emulator.Memory[emulator.PC+1])
		//fmt.Printf("jump to %x\n", foo)

		return
	}

	switch opcode & 0xF000 {
	case 0x4000:
		// 4XNN	Cond	if(Vx!=NN)	Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		v := (opcode & 0x0F00) >> 8
		if emulator.V[v] != byte(opcode) {

		}
	case 0x6000:
		// 6XNN	Const	Vx = NN	Sets VX to NN.
		v := (opcode & 0x0F00) >> 8
		emulator.V[v] = byte(opcode)
	case 0x7000:
		// 7XNN	Const	Vx += NN	Adds NN to VX.
		v := (opcode & 0x0F00) >> 8
		emulator.V[v] += byte(opcode)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0004:
			// 8XY4	Math	Vx += Vy	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			vx := (opcode & 0x0F00) >> 8
			vy := (opcode & 0x00F0) >> 4

			s := emulator.V[vx] + emulator.V[vy]
			emulator.V[vx] = byte(s)
			emulator.V[0xF] = (s >> 8) & 0x1
		default:
			fmt.Printf("%x\n", opcode)
		}
	case 0xA000:
		//ANNN	MEM	I = NNN	Sets I to the address NNN.
		emulator.I = 0x0FFF & opcode
	case 0xD000:
		x := emulator.V[(opcode&0x0F00)>>8]
		y := emulator.V[(opcode&0x00F0)>>4]
		height := opcode & 0x000F
		mem := emulator.Memory[emulator.I : emulator.I+height]

		emulator.Output.SetPixel(x, y, mem)
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x15:
			// FX07	Timer	Vx = get_delay()	Sets VX to the value of the delay timer.
			x := (opcode & 0x0F00) >> 8
			emulator.DTimer = emulator.V[x]
		default:
			fmt.Printf("%x\n", opcode)
		}

	default:
		fmt.Printf("%x\n", opcode)
	}
	emulator.PC += 2
}
