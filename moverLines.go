package main

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type MoverLines struct {
	speedx     float64
	accx       float64
	speedy     float64
	accy       float64
	currentAng int
	accAng     float64
	speedAng   float64
	centerx    float64
	centery    float64
	lines      []*canvas.Line
}

/*
-------------------------------------------------------------------- MoverLines
*/

func NewMoverLines(centerx, centery, speedAng float64) *MoverLines {
	return &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerx: centerx, centery: centery, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
}

func (mv *MoverLines) String() string {
	saw := mv.GetSizeAndCenter()
	return fmt.Sprintf("Lines x:%.3f y:%.3f w:%.3f h:%.3f", saw.CenterX, saw.CenterY, saw.Width, saw.Height)
}

func NewMoverRect(colour color.Color, centerx, centery, w, h, speedAng float64) *MoverLines {
	ml := &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerx: centerx, centery: centery, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
	ml.AddLine(colour, float32(centerx-w/2), float32(centery-h/2), float32(centerx+w/2), float32(centery-h/2))
	ml.AddLineToo(colour, float32(centerx+w/2), float32(centery+h/2))
	ml.AddLineToo(colour, float32(centerx-w/2), float32(centery+h/2))
	ml.AddLineToo(colour, float32(centerx-w/2), float32(centery-h/2))
	return ml
}

func (mv *MoverLines) ContainsAny(p *Points) bool {
	if mv.IsVisible() {
		return mv.GetBounds().ContainsAny(p)
	}
	return false
}

func (mv *MoverLines) SetVisible(v bool) {
	for _, l := range mv.lines {
		if !v {
			l.Hide()
		} else {
			l.Show()
		}
	}
}

func (mv *MoverLines) IsVisible() bool {
	return mv.lines[0].Visible()
}

func (mv *MoverLines) Update(time float64) {
	mv.accx = mv.accx + (mv.speedx * time)
	var dx int = 0
	if math.Abs(mv.accx) > 1 {
		dx = int(mv.accx)
		mv.accx = mv.accx - float64(dx)
	}
	mv.accy = mv.accy + (mv.speedy * time)
	var dy int = 0
	if math.Abs(mv.accy) > 1 {
		dy = int(mv.accy)
		mv.accy = mv.accy - float64(dy)
	}
	// Adjust the angle keeping it within bounds
	mv.accAng = mv.accAng + (mv.speedAng * float64(time))
	for mv.accAng > 360.0 {
		mv.accAng = mv.accAng - 360.0
	}
	for mv.accAng < 0 {
		mv.accAng = mv.accAng + 360.0
	}
	// Calc how much we need to rotate!
	intAng := int(mv.accAng)
	ra := 0
	if mv.currentAng != intAng {
		ra = intAng - mv.currentAng
		mv.currentAng = intAng
	}

	for _, l := range mv.lines {
		if dx != 0 {
			l.Position1.X = l.Position1.X + float32(dx)
			l.Position2.X = l.Position2.X + float32(dx)
		}
		if dy != 0 {
			l.Position1.Y = l.Position1.Y + float32(dy)
			l.Position2.Y = l.Position2.Y + float32(dy)
		}
		if ra != 0 {
			rotatePosition(mv.centerx, mv.centery, &l.Position1, ra)
			rotatePosition(mv.centerx, mv.centery, &l.Position2, ra)
		}
	}
	mv.centerx = mv.centerx + float64(dx)
	mv.centery = mv.centery + float64(dy)
}

func (mv *MoverLines) GetCanvasObject() fyne.CanvasObject {
	container := container.New(&ControllerLayout{})
	for _, l := range mv.lines {
		container.Add(l)
	}
	return container
}

func (mv *MoverLines) GetAngle() int {
	return mv.currentAng
}

func (mv *MoverLines) SetAngle(a int) {
	mv.currentAng = a
}

func (mv *MoverLines) SetAngleSpeed(as float64) {
	mv.speedAng = as
}

func (mv *MoverLines) GetAngleSpeed() float64 {
	return mv.speedAng
}

func (mv *MoverLines) SetSize(size fyne.Size) {
	currentSize := mv.GetSizeAndCenter()
	scaleX := float64(size.Width) / currentSize.Width
	scaleY := float64(size.Height) / currentSize.Height
	for _, l := range mv.lines {
		scalePoint(mv.centerx, mv.centery, &l.Position1, scaleX, scaleY)
		scalePoint(mv.centerx, mv.centery, &l.Position2, scaleX, scaleY)
	}
}

func (mv *MoverLines) GetSizeAndCenter() *SizeAndCenter {
	maxx := minFloat32
	maxy := minFloat32
	minx := maxFloat32
	miny := maxFloat32
	for _, l := range mv.lines {
		p1 := l.Position1
		p2 := l.Position2
		if p1.X > maxx {
			maxx = p1.X
		}
		if p1.X < minx {
			minx = p1.X
		}
		if p2.X > maxx {
			maxx = p2.X
		}
		if p2.X < minx {
			minx = p2.X
		}

		if p1.Y > maxy {
			maxy = p1.Y
		}
		if p1.Y < miny {
			miny = p1.Y
		}
		if p2.Y > maxy {
			maxy = p2.Y
		}
		if p2.Y < miny {
			miny = p2.Y
		}
	}
	w := float64(maxx - minx)
	h := float64(maxy - miny)
	return NewSizeAndCenter(w, h, (float64(minx) + w/2), (float64(miny) + h/2))
}

func (mv *MoverLines) GetBounds() *Bounds {
	sac := mv.GetSizeAndCenter()
	w2 := sac.Width / 2
	h2 := sac.Height / 2
	return &Bounds{x1: sac.CenterX - w2, x2: sac.CenterX + w2, y1: sac.CenterY - h2, y2: sac.CenterY + h2}
}

func (mv *MoverLines) GetPoints() *Points {
	if mv.lines[0].Visible() {
		n := len(mv.lines)
		p := &Points{x: make([]float64, n), y: make([]float64, n)}
		for i := 0; i < n; i++ {
			p.x[i] = float64(mv.lines[i].Position1.X)
			p.y[i] = float64(mv.lines[i].Position1.Y)
		}
		return p
	}
	return &Points{x: make([]float64, 0), y: make([]float64, 0)}
}

func (mv *MoverLines) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverLines) SetCenter(x, y float64) {
	dx := x - mv.centerx
	dy := y - mv.centery
	for _, l := range mv.lines {
		l.Position1.X = l.Position1.X + float32(dx)
		l.Position2.X = l.Position2.X + float32(dx)
		l.Position1.Y = l.Position1.Y + float32(dy)
		l.Position2.Y = l.Position2.Y + float32(dy)
	}
	mv.centerx = x
	mv.centery = y
}

func (mv *MoverLines) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
}

func (mv *MoverLines) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

func (mv *MoverLines) AddLine(colour color.Color, x1, y1, x2, y2 float32) {
	mv.lines = append(mv.lines, &canvas.Line{StrokeColor: colour, StrokeWidth: 1, Position1: fyne.Position{X: x1, Y: y1}, Position2: fyne.Position{X: x2, Y: y2}})
}

func (mv *MoverLines) AddLines(colour color.Color, xy ...float32) {
	for i := 0; i < len(xy); i = i + 2 {
		mv.AddLineToo(colour, xy[i], xy[i+1])
	}
}

func (mv *MoverLines) AddLineToo(colour color.Color, x2, y2 float32) {
	var x1 float32 = 0.0
	var y1 float32 = 0.0
	le := len(mv.lines)
	if le > 0 {
		lf := mv.lines[le-1]
		x1 = lf.Position2.X
		y1 = lf.Position2.Y
	}
	mv.AddLine(colour, x1, y1, x2, y2)
}
