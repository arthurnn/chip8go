package chip8

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

func (emulator *Chip8) Run() {
	//fmt.Printf("PC IS %x\n", emulator.PC)
	opcode := emulator.peekNextOp()
	emulator.cpu.Cycle(opcode)
	if emulator.cpu.ClearDisplay {
		emulator.Output.ClearDisplay()
	}

	if emulator.cpu.DrawSprit != nil {
		mem := emulator.Memory[emulator.cpu.I : emulator.cpu.I+emulator.cpu.DrawSprit.height]

		emulator.Output.SetPixel(emulator.cpu.DrawSprit.x, emulator.cpu.DrawSprit.y, mem)

//		if col {
			// V[0xF] = 1
		//}
	}
}
