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
