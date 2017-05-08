package chip8

import (
	"fmt"
)

type Chip8 struct {
	cpu CPU

	Memory [4096]byte

	//	TODO keyboard input

	Output Graphics
}

func NewChip8() *Chip8 {
	cpu := CPU{
		PC: 0x200,
	}

	e := &Chip8{
		cpu: cpu,
	}
	return e
}

func (emulator *Chip8) LoadRom(rom []byte) {
	memArea := emulator.Memory[emulator.cpu.PC:]
	copy(memArea, rom)
}

func (emulator *Chip8) peekNextOp() uint16 {
	return uint16(emulator.Memory[emulator.cpu.PC])<<8 | uint16(emulator.Memory[emulator.cpu.PC+1])
}

func (emulator *Chip8) Cycle() {
	//fmt.Printf("PC IS %x\n", emulator.PC)
	opcode := emulator.peekNextOp()
	if 0x1000 == (opcode & 0xF000) {
		// 1NNN	Flow	goto NNN;	Jumps to address NNN.
		emulator.cpu.PC = opcode & 0x0FFF
		return
	}

	switch opcode & 0xF000 {
	case 0x4000:
		// 4XNN
		// Cond	if(Vx!=NN)
		// Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		v := (opcode & 0x0F00) >> 8
		if emulator.cpu.V[v] != byte(opcode) {
			// TODO
		}
	case 0x6000:
		// 6XNN	Const	Vx = NN	Sets VX to NN.
		v := (opcode & 0x0F00) >> 8
		emulator.cpu.V[v] = byte(opcode)
	case 0x7000:
		// 7XNN	Const	Vx += NN	Adds NN to VX.
		v := (opcode & 0x0F00) >> 8
		emulator.cpu.V[v] += byte(opcode)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0004:
			// 8XY4	Math	Vx += Vy	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			vx := (opcode & 0x0F00) >> 8
			vy := (opcode & 0x00F0) >> 4

			s := emulator.cpu.V[vx] + emulator.cpu.V[vy]
			emulator.cpu.V[vx] = byte(s)
			emulator.cpu.V[0xF] = (s >> 8) & 0x1
		default:
			fmt.Printf("%x\n", opcode)
		}
	case 0xA000:
		//ANNN	MEM	I = NNN	Sets I to the address NNN.
		emulator.cpu.I = 0x0FFF & opcode
	case 0xD000:
		x := emulator.cpu.V[(opcode&0x0F00)>>8]
		y := emulator.cpu.V[(opcode&0x00F0)>>4]
		height := opcode & 0x000F
		mem := emulator.Memory[emulator.cpu.I : emulator.cpu.I+height]

		emulator.Output.SetPixel(x, y, mem)
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x15:
			// FX07	Timer	Vx = get_delay()	Sets VX to the value of the delay timer.
			x := (opcode & 0x0F00) >> 8
			emulator.cpu.DTimer = emulator.cpu.V[x]
		default:
			fmt.Printf("%x\n", opcode)
		}

	default:
		fmt.Printf("%x\n", opcode)
	}

	emulator.cpu.PC += 2

	if emulator.cpu.DTimer > 0 {
		emulator.cpu.DTimer -= 1
	}
}
