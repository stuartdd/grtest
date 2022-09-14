package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type MoverCircle struct {
	speedx     float64
	speedy     float64
	centerx    float64
	centery    float64
	width      float64
	height     float64
	circle     *canvas.Circle
	shouldMove func(Movable, float64, float64) bool
}

/*
-------------------------------------------------------------------- MoverCircle
*/
func NewMoverCircle(strokeColour, fillColour color.Color, centerx, centery, width, height float64) *MoverCircle {
	circle := &canvas.Circle{
		Position1:   fyne.Position{X: float32(centerx - width/2), Y: float32(centery - height/2)},
		Position2:   fyne.Position{X: float32(centerx + width/2), Y: float32(centery + height/2)},
		StrokeColor: strokeColour, FillColor: fillColour, StrokeWidth: 1,
	}
	return &MoverCircle{speedx: 0, speedy: 0, centerx: centerx, centery: centery, width: width, height: height, circle: circle}
}

func (mv *MoverCircle) ContainsAny(p *Points) bool {
	if mv.IsVisible() {
		return mv.GetBounds().ContainsAny(p)
	}
	return false
}

func (mv *MoverCircle) String() string {
	return fmt.Sprintf("Circle x:%.3f y:%.3f w:%.3f h:%.3f", mv.centerx, mv.centery, mv.width, mv.height)
}

func (mv *MoverCircle) SetVisible(v bool) {
	if !v {
		mv.circle.Hide()
	} else {
		mv.circle.Show()
	}
}

func (mv *MoverCircle) IsVisible() bool {
	return mv.circle.Visible()
}

func (mv *MoverCircle) SetShouldMove(f func(Movable, float64, float64) bool) {
	mv.shouldMove = f
}

func (mv *MoverCircle) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if mv.shouldMove == nil || (mv.shouldMove != nil && mv.shouldMove(mv, dx, dy)) {
		mv.centerx = mv.centerx + dx
		mv.centery = mv.centery + dy
		mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
		mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
		mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
		mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
	}
}

func (mv *MoverCircle) UpdateContainerWithObjects(c *fyne.Container) {
	c.Add(mv.circle)
}

func (mv *MoverCircle) GetSizeAndCenter() *SizeAndCenter {
	return NewSizeAndCenter(mv.width, mv.height, mv.centerx, mv.centery)
}

func (mv *MoverCircle) GetBounds() *Bounds {
	w2 := mv.width / 2
	h2 := mv.height / 2
	return &Bounds{x1: mv.centerx - w2, x2: mv.centerx + w2, y1: mv.centery - h2, y2: mv.centery + h2}
}

func (mv *MoverCircle) GetPoints() *Points {
	if mv.circle.Visible() {
		p := &Points{x: make([]float64, 4), y: make([]float64, 4)}
		cx := mv.centerx
		cy := mv.centery
		w2 := mv.width / 2
		h2 := mv.height / 2
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

func (mv *MoverCircle) SetSize(size fyne.Size) {
	mv.width = float64(size.Width)
	mv.height = float64(size.Height)
	mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
	mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
	mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
	mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
}

func (mv *MoverCircle) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
}

func (mv *MoverCircle) SetCenter(x, y float64) {
	mv.centerx = x
	mv.centery = y
	mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
	mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
	mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
	mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
}

func (mv *MoverCircle) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverCircle) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

func (mv *MoverCircle) SetAngle(a int) {
}

func (mv *MoverCircle) GetAngle() int {
	return 0
}

func (mv *MoverCircle) SetAngleSpeed(as float64) {
}

func (mv *MoverCircle) GetAngleSpeed() float64 {
	return 0
}
