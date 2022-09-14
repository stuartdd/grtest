package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
)

var (
	bounds = NewBounds(20, 20, 480, 480)
)

/*
-------------------------------------------------------------------- main
*/
func mainPOClots(mainWindow fyne.Window, controller *ControllerLayout) {
	for i := 0; i < 2000; i++ {
		x := rand.Float64() * 400
		y := rand.Float64() * 400
		c := NewMoverCircle(nil, color.RGBA{255, 0, 0, 255}, 50+x, 50+y, 5, 5)
		c.SetSpeed(rnd(), rnd())
		controller.AddMover(c)
		c.SetShouldMove(shouldMove)
	}

	go func() {
		for {
			w := float64(mainWindow.Canvas().Size().Width)
			if w < 100 {
				w = 100
			}
			h := float64(mainWindow.Canvas().Size().Height)
			if h < 100 {
				h = 100
			}
			b := NewBounds(20, 20, w-40, h-40)
			if !bounds.Equal(b) {
				bounds = b
			}
			time.Sleep(time.Second)
		}
	}()

}

func rnd() float64 {
	if rand.Float64() > 0.5 {
		return -(50 + rand.Float64()*50)
	}
	return (50 + rand.Float64()*50)
}

func shouldMove(m Movable, dx, dy float64) bool {
	sx, sy := m.GetSpeed()
	px, py := m.GetCenter()
	b := bounds.Outside(px+dx, py+dy)
	switch b {
	case LEFT, RIGHT:
		m.SetSpeed(-sx, sy)
		return false
	case UP, DOWN:
		m.SetSpeed(sx, -sy)
		return false
	case INSIDE:
		return true
	default:
		m.(*MoverCircle).circle.FillColor = color.RGBA{0, 255, 0, 255}
		m.SetSpeed(-sx, -sy)
		return false
	}
}
