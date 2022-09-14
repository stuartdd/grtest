package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type MoverText struct {
	speedx     float64
	speedy     float64
	positionx  float64
	centery    float64
	width      float64
	height     float64
	text       *canvas.Text
	align      fyne.TextAlign
	shouldMove func(Movable, float64, float64) bool
}

/*
-------------------------------------------------------------------- MoverText
*/
func NewMoverText(text string, posx, centery float64, si float32, align fyne.TextAlign) *MoverText {
	t := &canvas.Text{Text: text, TextSize: si, TextStyle: textStyle}
	mv := &MoverText{text: t, align: align, positionx: posx, centery: centery}
	mv.textSize()
	mv.moveToPosition()
	return mv
}

func (mv *MoverText) String() string {
	return fmt.Sprintf("Text t:%s pos:%.3f y:%.3f w:%.3f h:%.3f", mv.text.Text, mv.positionx, mv.centery, mv.width, mv.height)
}

func (mv *MoverText) SetShouldMove(f func(Movable, float64, float64) bool) {
	mv.shouldMove = f
}

func (mv *MoverText) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if mv.shouldMove == nil || (mv.shouldMove != nil && mv.shouldMove(mv, dx, dy)) {
		mv.positionx = mv.positionx + dx
		mv.centery = mv.centery + dy
		mv.moveToPosition()
	}
}

func (mv *MoverText) ContainsAny(p *Points) bool {
	if mv.IsVisible() {
		return mv.GetBounds().ContainsAny(p)
	}
	return false
}

func (mv *MoverText) SetVisible(v bool) {
	if !v {
		mv.text.Hide()
	} else {
		mv.text.Show()
	}
}

func (mv *MoverText) IsVisible() bool {
	return mv.text.Visible()
}

func (mv *MoverText) UpdateContainerWithObjects(cont *fyne.Container) {
	cont.Add(mv.text)
}

func (mv *MoverText) GetSizeAndCenter() *SizeAndCenter {
	cx, cy := mv.textCenter()
	return NewSizeAndCenter(float64(mv.width), float64(mv.height), cx, cy)
}

func (mv *MoverText) GetBounds() *Bounds {
	w2 := float64(mv.width / 2)
	h2 := float64(mv.height / 2)
	cx, cy := mv.textCenter()
	return &Bounds{x1: cx - w2, x2: cx + w2, y1: cy - h2, y2: cy + h2}
}

func (mv *MoverText) GetPoints() *Points {
	if mv.text.Visible() {
		p := &Points{x: make([]float64, 4), y: make([]float64, 4)}
		cx, cy := mv.textCenter()
		w2 := float64(mv.width / 2)
		h2 := float64(mv.height / 2)
		p.x[0] = cx - w2 // Top left
		p.y[0] = cy - h2
		p.x[1] = cx + w2 // Top Right
		p.y[1] = cy - h2
		p.x[2] = cx + w2 // Bottom Right
		p.y[2] = cy + h2
		p.x[3] = cx - w2 // Bottom Left
		p.y[3] = cy + h2
		return p
	}
	return &Points{x: make([]float64, 0), y: make([]float64, 0)}
}

func (mv *MoverText) SetSize(s fyne.Size) {
	mv.text.Resize(s)
	mv.textSize()
	mv.moveToPosition()
}

func (mv *MoverText) SetText(text string) {
	mv.text.Text = text
	mv.textSize()
	mv.moveToPosition()
}

func (mv *MoverText) GetCenter() (float64, float64) {
	return mv.positionx, mv.centery
}

func (mv *MoverText) SetCenter(x, y float64) {
	mv.positionx = x
	mv.centery = y
	mv.moveToPosition()
}

func (mv *MoverText) SetAngle(a int) {
}

func (mv *MoverText) GetAngle() int {
	return 0
}

func (mv *MoverText) SetAngleSpeed(as float64) {
}

func (mv *MoverText) GetAngleSpeed() float64 {
	return 0
}

func (mv *MoverText) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverText) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

func (mv *MoverText) textCenter() (float64, float64) {
	switch mv.align {
	case fyne.TextAlignLeading:
		return mv.positionx + mv.width/2, mv.centery
	case fyne.TextAlignTrailing:
		return mv.positionx - mv.width/2, mv.centery
	}
	return mv.positionx, mv.centery
}

func (mv *MoverText) textPosition() (float64, float64) {
	switch mv.align {
	case fyne.TextAlignCenter:
		return mv.positionx - mv.width/2, mv.centery - mv.height/2
	case fyne.TextAlignTrailing:
		return mv.positionx - mv.width, mv.centery - mv.height/2
	}
	return mv.positionx, mv.centery - mv.height/2
}

func (mv *MoverText) textSize() {
	size := fyne.MeasureText(mv.text.Text, mv.text.TextSize, mv.text.TextStyle)
	mv.width = float64(size.Width)
	mv.height = float64(size.Height)
}

func (mv *MoverText) moveToPosition() {
	px, py := mv.textPosition()
	mv.text.Move(fyne.Position{X: float32(px), Y: float32(py)})
}
