package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Quad int

const (
	LEFT Quad = iota
	RIGHT
	UP
	DOWN
	LEFT_UP
	LEFT_DOWN
	RIGHT_UP
	RIGHT_DOWN
	INSIDE
)

type SizeAndCenter struct {
	Width   float64
	Height  float64
	CenterX float64
	CenterY float64
}

type Bounds struct {
	x1 float64
	y1 float64
	x2 float64
	y2 float64
}

type Points struct {
	x []float64
	y []float64
}

func QuadString(q Quad) string {
	switch q {
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case LEFT_DOWN:
		return "LEFT+DOWN"
	case LEFT_UP:
		return "LEFT+UP"
	case RIGHT_DOWN:
		return "RIGHT+DOWN"
	case RIGHT_UP:
		return "RIGHT+UP"
	case INSIDE:
		return "INSIDE"
	}
	return "UNKNOWN"
}

func MeasureString(text string) fyne.Size {
	ts := fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2}
	si := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)
	return fyne.MeasureText(text, si, ts)
}

/*
-------------------------------------------------------------------- SizeAndCenter + Bounds
*/
func NewSizeAndCenter(w, h, x, y float64) *SizeAndCenter {
	return &SizeAndCenter{Width: w, Height: h, CenterX: x, CenterY: y}
}

func (sc *SizeAndCenter) Size() fyne.Size {
	return fyne.Size{Height: float32(sc.Height), Width: float32(sc.Width)}
}

func (sc *SizeAndCenter) Center() fyne.Position {
	return fyne.Position{X: float32(sc.CenterX), Y: float32(sc.CenterY)}
}

func (sc *SizeAndCenter) Points() *Points {
	w2 := sc.Width / 2
	h2 := sc.Height / 2
	p := &Points{x: make([]float64, 4), y: make([]float64, 4)}
	p.x[0] = sc.CenterX - w2 // Top left
	p.y[0] = sc.CenterY - h2
	p.x[1] = sc.CenterX + w2 // Top Right
	p.y[1] = sc.CenterY - h2
	p.x[2] = sc.CenterX + w2 // Bottom Right
	p.y[2] = sc.CenterY + h2
	p.x[3] = sc.CenterX - w2 // Bottom Left
	p.y[3] = sc.CenterY + h2
	return p
}

func NewBounds(x1, y1, x2, y2 float64) *Bounds {
	return &Bounds{x1: x1, y1: y1, x2: x2, y2: y2}
}

func (bb *Bounds) Size() fyne.Size {
	return fyne.Size{Height: float32(bb.y2 - bb.y1), Width: float32(bb.x2 - bb.x1)}
}

func (bb *Bounds) Center() fyne.Position {
	return fyne.Position{X: float32(bb.x1 + (bb.x2-bb.x1)/2), Y: float32(bb.y1 + (bb.y2-bb.y1)/2)}
}

func (bb *Bounds) Equal(aa *Bounds) bool {
	return (bb.x1 == aa.x1 && bb.x2 == aa.x2 && bb.y1 == aa.y1 && bb.y2 == aa.y2)
}

func (bb *Bounds) Points() *Points {
	p := &Points{x: make([]float64, 4), y: make([]float64, 4)}
	p.x[0] = bb.x1 // Top left
	p.y[0] = bb.y1
	p.x[1] = bb.x2 // Top Right
	p.y[1] = bb.y1
	p.x[2] = bb.x2 // Bottom Right
	p.y[2] = bb.y2
	p.x[3] = bb.x1 // Bottom Left
	p.y[3] = bb.y2
	return p
}

func (bb *Bounds) Inside(x, y float64) bool {
	/*
		Using && is faster as only 1 expressions (may) need to be evaluated to return false
		Using || ALL expressions need to be evaluated
	*/
	if x >= bb.x1 && x <= bb.x2 && y >= bb.y1 && y <= bb.y2 {
		return true
	}
	return false
}

func (bb *Bounds) Outside(x, y float64) Quad {
	if x < bb.x1 {
		// Left
		if y < bb.y1 {
			// left up
			return LEFT_UP
		}
		if y > bb.y2 {
			// left down
			return LEFT_DOWN
		}
		return LEFT
	}
	if x > bb.x2 {
		// right
		if y < bb.y1 {
			// right up
			return RIGHT_UP
		}
		if y > bb.y2 {
			// right down
			return RIGHT_DOWN
		}
		return RIGHT
	}
	if y < bb.y1 {
		// right up
		return UP
	}
	if y > bb.y2 {
		// right down
		return DOWN
	}
	return INSIDE
}

func (bb *Bounds) ContainsAny(p *Points) bool {
	l := len(p.x)
	for i := 0; i < l; i++ {
		if bb.Inside(p.x[i], p.y[i]) {
			return true
		}
	}
	return false
}

func (bb *SizeAndCenter) String() string {
	return fmt.Sprintf("SizeAndCenter w:%.3f h:%.3f x:%.3f y:%.3f", bb.Width, bb.Height, bb.CenterX, bb.CenterY)
}

func (bb *Points) String() string {
	var sb strings.Builder
	sb.WriteString("Points:")
	l := len(bb.x)
	for i := 0; i < l; i++ {
		sb.WriteString(fmt.Sprintf(" (x:%.3f y:%.3f)", bb.x[i], bb.y[i]))
	}
	return sb.String()
}

func (bb *Bounds) String() string {
	return fmt.Sprintf("Bounds x1:%.3f y1:%.3f x2:%.3f y2:%.3f", bb.x1, bb.y1, bb.x2, bb.y2)
}
