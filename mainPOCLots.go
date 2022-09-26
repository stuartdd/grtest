package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var (
	bounds = NewBounds(20, 20, 480, 480)
)

/*
-------------------------------------------------------------------- main
*/
func MainPOCLots(mainWindow fyne.Window, width, height float64, controller *MoverController) *fyne.Container {

	cont := container.New(NewStaticLayout(width, height))

	for i := 0; i < 2000; i++ {
		c := NewMoverCircle(nil, color.RGBA{255, 0, 0, 255}, rndPos(35, height-40), rndPos(35, height-40), 5, 5)
		c.SetSpeed(rndSpeed(10, 50), rndSpeed(10, 50))
		controller.AddMover(c)
		for _, co := range c.GetCanvasObjects() {
			cont.Add(co)
		}
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
	return cont
}

func rndSpeed(min, max float64) float64 {
	if rand.Float64() > 0.5 {
		return -(min + rand.Float64()*(max-min))
	}
	return min + rand.Float64()*(max-min)
}

func rndPos(min, max float64) float64 {
	return min + (rand.Float64() * (max - min))
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
