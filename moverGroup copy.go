package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MoverGroup struct {
	movers     []Movable
	currentAng int
	accAng     float64
	speedAng   float64
}

/*
-------------------------------------------------------------------- MoverGroup
*/
func NewMoverCroup(mandatoryMover Movable) *MoverGroup {
	mv := &MoverGroup{movers: make([]Movable, 0), currentAng: 0, speedAng: 0}
	mv.Add(mandatoryMover)
	return mv
}

func (mv *MoverGroup) ContainsAny(p *Points) bool {
	if mv.movers[0].IsVisible() {
		return mv.movers[0].GetBounds().ContainsAny(p)
	}
	return false
}

func (mv *MoverGroup) SetVisible(v bool) {
	for _, m := range mv.movers {
		m.SetVisible(v)
	}
}

func (mv *MoverGroup) IsVisible() bool {
	return mv.movers[0].IsVisible()
}

func (mv *MoverGroup) Add(mover Movable) {
	mv.movers = append(mv.movers, mover)
	mover.SetSpeed(mv.movers[0].GetSpeed())
}

func (mv *MoverGroup) Update(time float64) {
	mv0 := mv.movers[0]
	mv0.Update(time)

	cx, cy := mv0.GetCenter()
	ra := 0
	if mv.speedAng != 0 {
		mv.accAng = mv.accAng + (mv.speedAng * float64(time))
		for mv.accAng > 360.0 {
			mv.accAng = mv.accAng - 360.0
		}
		for mv.accAng < 0 {
			mv.accAng = mv.accAng + 360.0
		}
		// Calc how much we need to rotate!
		intAng := int(mv.accAng)
		if mv.currentAng != intAng {
			ra = intAng - mv.currentAng
			mv.currentAng = intAng
		}
	}
	for i := 1; i < len(mv.movers); i++ {
		m := mv.movers[i]
		m.Update(time)
		if ra != 0 {
			cxm, cym := m.GetCenter()
			m.SetCenter(rotatePoints(cx, cy, cxm, cym, ra))
		}
	}
}

func (mv *MoverGroup) GetCanvasObject() fyne.CanvasObject {
	container := container.New(&ControllerLayout{})
	for _, m := range mv.movers {
		container.Add(m.GetCanvasObject())
	}
	return container
}

func (mv *MoverGroup) GetSizeAndCenter() *SizeAndCenter {
	return mv.movers[0].GetSizeAndCenter()
}

func (mv *MoverGroup) SetCenter(x, y float64) {
	cx, cy := mv.movers[0].GetCenter()
	dx := x - cx
	dy := y - cy
	for _, m := range mv.movers {
		ccx, ccy := m.GetCenter()
		m.SetCenter(ccx+dx, ccy+dy)
	}
}

func (mv *MoverGroup) GetCenter() (float64, float64) {
	return mv.movers[0].GetCenter()
}

func (mv *MoverGroup) SetSpeed(sx float64, sy float64) {
	for _, m := range mv.movers {
		m.SetSpeed(sx, sy)
	}
}

func (mv *MoverGroup) GetSpeed() (float64, float64) {
	return mv.movers[0].GetSpeed()
}

func (mv *MoverGroup) SetAngleSpeed(as float64) {
	mv.speedAng = as
}

func (mv *MoverGroup) GetAngleSpeed() float64 {
	return mv.speedAng
}

func (mv *MoverGroup) SetAngle(a int) {
	mv.currentAng = a
	for _, m := range mv.movers {
		m.SetAngle(a)
	}
}

func (mv *MoverGroup) GetAngle() int {
	return mv.currentAng
}

func (mv *MoverGroup) GetBounds() *Bounds {
	return mv.movers[0].GetBounds()
}

func (mv *MoverGroup) GetPoints() *Points {
	return mv.movers[0].GetPoints()
}

func (mv *MoverGroup) SetSize(fyne.Size) {

}
