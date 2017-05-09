package chip8

import (
	"fmt"
)

type CPU struct {
	// registers
	V  [16]byte
	PC uint16
	I  uint16

	// Drawing
	DrawSprit *Sprit
	ClearDisplay bool

	// stack
	sp    uint16
	stack [16]uint16

	// timers
	DTimer byte
	STime  byte
}

func (cpu *CPU) pushToStack(pc uint16) {
	cpu.stack[cpu.sp] = pc
	cpu.sp += 1

	// TODO check stack > 15
}

func (cpu *CPU) popFromStack() uint16 {
	cpu.sp -= 1
	return cpu.stack[cpu.sp]
}

func (cpu *CPU) Cycle(opcode uint16) {
	cpu.DrawSprit = nil
	cpu.ClearDisplay = false

	if 0x00EE == opcode {
		// 00EE
		// Flow
		// return;
		// Returns from a subroutine.
		cpu.PC = cpu.popFromStack()
		return
	}

	if 0x1000 == (opcode & 0xF000) {
		// 1NNN
		// Flow
		// goto NNN;
		// Jumps to address NNN.
		cpu.PC = opcode & 0x0FFF
		return
	}

	if 0x2000 == (opcode & 0xF000) {
		// 2NNN
		// Flow
		// *(0xNNN)()
		// Calls subroutine at NNN.
		cpu.pushToStack(cpu.PC)
		cpu.PC = opcode & 0x0FFF
		return
	}

	switch opcode & 0xF000 {
	case 0x0000:
		if opcode == 0x00E0 {
			// 00E0
			// Display
			// Clears the screen.
			cpu.ClearDisplay = true
		}
	case 0x4000:
		// 4XNN
		// Cond	if(Vx!=NN)
		// Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		v := (opcode & 0x0F00) >> 8
		if cpu.V[v] != byte(opcode) {
			// TODO
		}
	case 0x6000:
		// 6XNN	Const	Vx = NN	Sets VX to NN.
		v := (opcode & 0x0F00) >> 8
		cpu.V[v] = byte(opcode)
	case 0x7000:
		// 7XNN	Const	Vx += NN	Adds NN to VX.
		v := (opcode & 0x0F00) >> 8
		cpu.V[v] += byte(opcode)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0004:
			// 8XY4	Math	Vx += Vy	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			vx := (opcode & 0x0F00) >> 8
			vy := (opcode & 0x00F0) >> 4

			s := cpu.V[vx] + cpu.V[vy]
			cpu.V[vx] = byte(s)
			cpu.V[0xF] = (s >> 8) & 0x1
		default:
			fmt.Printf("%x\n", opcode)
		}
	case 0xA000:
		//ANNN	MEM	I = NNN	Sets I to the address NNN.
		cpu.I = 0x0FFF & opcode
	case 0xD000:
		// Draw
		x := cpu.V[(opcode&0x0F00)>>8]
		y := cpu.V[(opcode&0x00F0)>>4]
		height := opcode & 0x000F

		cpu.DrawSprit = &Sprit{
			x:      x,
			y:      y,
			height: height,
		}

		// TODO
		//mem := emulator.Memory[cpu.I : cpu.I+height]
		//emulator.Output.SetPixel(x, y, mem)
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x15:
			// FX07	Timer	Vx = get_delay()	Sets VX to the value of the delay timer.
			x := (opcode & 0x0F00) >> 8
			cpu.DTimer = cpu.V[x]
		default:
			fmt.Printf("%x\n", opcode)
		}

	default:
		fmt.Printf("Not found instruction %x\n", opcode)
	}

	cpu.PC += 2

	if cpu.DTimer > 0 {
		cpu.DTimer -= 1
	}
}
