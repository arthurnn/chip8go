package chip8

type Memory [4096]byte

type Chip8 struct {
	cpu CPU

	Memory Memory

	//	TODO keyboard input

	Output *Graphics
}

var defaultsprite = []byte{
	0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
	0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
	0x90, 0x90, 0xf0, 0x10, 0x10, // 4
	0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
	0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
	0xf0, 0x10, 0x20, 0x40, 0x40, // 7
	0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
	0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
	0xf0, 0x90, 0xf0, 0x90, 0x90, // A
	0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
	0xf0, 0x80, 0x80, 0x80, 0xf0, // C
	0xe0, 0x90, 0x90, 0x90, 0xe0, // D
	0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
	0xf0, 0x80, 0xf0, 0x80, 0x80, // F
}

func NewChip8() *Chip8 {
	cpu := CPU{
		PC: 0x200,
	}

	e := &Chip8{
		cpu:    cpu,
		Output: NewNullDisplay(),
	}

	// Load default sprites
	copy(e.Memory[:0x200], defaultsprite)
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
	emulator.cpu.Cycle(opcode, &emulator.Memory)

	if emulator.cpu.ClearDisplay {
		emulator.Output.ClearDisplay()
	}

	if emulator.cpu.DrawSprite != nil {
		mem := emulator.Memory[emulator.cpu.I : emulator.cpu.I+emulator.cpu.DrawSprite.height]

		col := emulator.Output.SetPixel(emulator.cpu.DrawSprite.x, emulator.cpu.DrawSprite.y, mem)
		if col {
			emulator.cpu.V[0xF] = 0x01
		} else {
			emulator.cpu.V[0xF] = 0x00
		}
	}
}
