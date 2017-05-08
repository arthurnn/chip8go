package chip8

type CPU struct {
	// registers
	V  [16]byte
	PC uint16
	I  uint16

	// stack
	sp    uint16
	stack [16]uint16

	// timers
	DTimer byte
	STime  byte
}

func (cpu *CPU) PushToStack(pc uint16) {
	cpu.stack[cpu.sp] = pc
	cpu.sp += 1

	// TODO check stack > 15
}

func (cpu *CPU) PopFromStack() uint16 {
	cpu.sp -= 1
	return cpu.stack[cpu.sp]
}
