package chip8

import (
	"fmt"
	"math/rand"
)

type CPU struct {
	// registers
	V  [16]byte
	PC uint16
	I  uint16

	// Drawing
	DrawSprite    *Sprite
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

	if cpu.sp > 15 {
		panic("Stackoverflow")
	}

}

func (cpu *CPU) popFromStack() uint16 {
	cpu.sp -= 1
	return cpu.stack[cpu.sp]
}

func (cpu *CPU) Cycle(opcode uint16, mem *Memory) {
	cpu.DrawSprite = nil
	cpu.ClearDisplay = false

	//fmt.Printf("Cycle on %x\n", opcode)

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
		switch opcode & 0x00FF {
		case 0xEE:
			// 00EE
			// Flow
			// return;
			// Returns from a subroutine.
			cpu.PC = cpu.popFromStack()
		case 0xE0:
			// 00E0
			// Display
			// Clears the screen.
			cpu.ClearDisplay = true
		default:
			panic(fmt.Sprintf("Not found instruction %x\n", opcode))
		}
	case 0x3000:
		// 3XNN	Cond	if(Vx==NN)	Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
		x := (opcode & 0x0F00) >> 8
		if x == (opcode & 0x00FF) {
			cpu.PC += 2
		}
	case 0x4000:
		// 4XNN	Cond	if(Vx!=NN)	Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		x := (opcode & 0x0F00) >> 8
		if x != (opcode & 0x00FF) {
			cpu.PC += 2
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
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		switch opcode & 0x000F {
		case 0x0002:
			// 8XY2	BitOp	Vx=Vx&Vy	Sets VX to VX and VY. (Bitwise AND operation)
			cpu.V[x] = cpu.V[x] & cpu.V[y]
		case 0x0004:
			// 8XY4	Math	Vx += Vy	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			s := cpu.V[x] + cpu.V[y]
			cpu.V[x] = byte(s)
			cpu.V[0xF] = (s >> 8) & 0x1
		default:
			panic(fmt.Sprintf("Not found instruction %x\n", opcode))
		}
	case 0xA000:
		//ANNN	MEM	I = NNN	Sets I to the address NNN.
		cpu.I = 0x0FFF & opcode
	case 0xC000:
		//CXNN	Rand	Vx=rand()&NN	Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		x := (opcode & 0x0F00) >> 8
		cpu.V[x] = byte(rand.Int()) & byte(opcode)
	case 0xD000:
		// Draw
		x := cpu.V[(opcode&0x0F00)>>8]
		y := cpu.V[(opcode&0x00F0)>>4]
		height := opcode & 0x000F

		cpu.DrawSprite = &Sprite{
			x:      x,
			y:      y,
			height: height,
		}
	case 0xE000:
		//x := (opcode & 0x0F00) >> 8
		switch opcode & 0x00FF {
		case 0xA1:
			//return TODO
		}
	case 0xF000:
		x := (opcode & 0x0F00) >> 8

		switch opcode & 0x00FF {
		case 0x07:
			// FX07	Timer	Vx = get_delay()	Sets VX to the value of the delay timer.
			cpu.V[x] = cpu.DTimer
		case 0x15:
			// FX15	Timer	delay_timer(Vx)	Sets the delay timer to VX.
			cpu.DTimer = cpu.V[x]
		case 0x29:
			// FX29	MEM	I=sprite_addr[Vx]	Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
			cpu.I = uint16(cpu.V[x]) * uint16(5)
		case 0x33:
			vx := cpu.V[x]
			mem[cpu.I] = byte(vx / 100)
			mem[cpu.I+1] = byte((vx / 10) % 10)
			mem[cpu.I+2] = byte(vx % 10)
			//panic("testttt")
		case 0x55:
			// FX55	MEM	reg_dump(Vx,&I)	Stores V0 to VX (including VX) in memory starting at address I.
			for i := uint16(0); i <= x; i += 1 {
				mem[cpu.I+i] = cpu.V[i]
			}
		case 0x65:
			// FX65	MEM	reg_load(Vx,&I)	Fills V0 to VX (including VX) with values from memory starting at address I.
			for i := uint16(0); i <= x; i += 1 {
				cpu.V[i] = mem[cpu.I+i]
			}
		default:
			panic(fmt.Sprintf("Not found instruction %x\n", opcode))
		}

	default:
		panic(fmt.Sprintf("Not found instruction %x\n", opcode))
	}

	cpu.PC += 2

	if cpu.DTimer > 0 {
		cpu.DTimer -= 1
	}
}
