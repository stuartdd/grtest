package main

import (
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
)

/*
-------------------------------------------------------------------- main
*/
func mainPOClots(mainWindow fyne.Window, controller *ControllerLayout) {

	bounds := NewBounds(20, 20, 480, 480)
	for i := 0; i < 100; i++ {
		x := rand.Float64() * 400
		y := rand.Float64() * 400
		c := NewMoverCircle(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, 50+x, 50+y, 10, 10)
		c.SetSpeed(rand.Float64()*100, rand.Float64()*100)
		controller.AddMover(c)
	}

	controller.SetOnUpdate(func(f float64) bool {
		for _, j := range controller.movers {
			sx, sy := j.GetSpeed()
			px, py := j.GetCenter()
			switch bounds.Outside(px, py) {
			case LEFT | RIGHT:
				j.SetSpeed(-sx, sy)
			case UP | DOWN:
				j.SetSpeed(sx, -sy)
			case INSIDE:
			default:
				j.SetSpeed(0, 0)
			}

		}
		return true
	})
}
