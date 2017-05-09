package chip8

import (
	"fmt"
)

type Sprit struct {
	x, y byte
	height uint16
}

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


func (g *Graphics) ClearDisplay() {

}
