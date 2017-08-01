package chip8

import (
//	"fmt"

	"github.com/nsf/termbox-go"
)

type Sprite struct {
	x, y   byte
	height uint16
}

type Graphics struct {
	drawFlag bool
	gfx      [65][33]bool

	display Display
}

type Display interface {
	Render(g *Graphics)
}

func termboxInit(bg termbox.Attribute) error {
	if err := termbox.Init(); err != nil {
		return err
	}

	termbox.HideCursor()

	if err := termbox.Clear(bg, bg); err != nil {
		return err
	}

	return termbox.Flush()
}

func NewTermboxDisplay() *Graphics {
	termboxInit(termbox.ColorDefault)
	return &Graphics{display: TermboxDisplay{}}
}

func NewNullDisplay() *Graphics {
	return &Graphics{display: NullDisplay{}}
}

type NullDisplay struct{}

func (n NullDisplay) Render(g *Graphics) {
	// empty
}

type TermboxDisplay struct{}

func (t TermboxDisplay) Render(g *Graphics) {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			var v rune
			if g.gfx[x][y] {
				v = '*'
			} else {
				v = ' '
			}
			termbox.SetCell(x, y, v, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	termbox.Flush()

}

func (g *Graphics) Render() {
	if g.drawFlag {
		g.drawFlag = false
		g.display.Render(g)
	}
}

func (g *Graphics) SetPixel(x, y byte, memory []byte) (collision bool) {
	g.drawFlag = true

	width, height := uint(8), uint(len(memory))

	//fmt.Printf("set pixel size %x\n", memory)
	for ry := uint(0); ry < height; ry++ {
		pixel := memory[ry]
		for rx := uint(0); rx < width; rx++ {
			p := 128 >> rx // 128 == 10000000 (binary)
			if (pixel & byte(p)) > 0 {
				gx := byte(rx) + x
				gy := byte(ry) + y
				//fmt.Printf("Set Pixel [%d,%d]\n", gx, gy)

				if g.gfx[gx][gy] {
					g.gfx[gx][gy] = false
					collision = true
				} else {
					g.gfx[gx][gy] = true
				}

			}
		}
	}

	return
}

func (g *Graphics) ClearDisplay() {
	//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	panic("clear")
}

func (g *Graphics) Close() {
	termbox.Close()
}
