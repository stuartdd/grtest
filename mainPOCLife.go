package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var (
	dots     []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos  int              = 0
	gridSize float32          = 5
	lifeGen  *LifeGen         = NewLifeGen()
	xOffset  float32          = 100
	yOffset  float32          = 100
	genColor []color.Color    = []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
)

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, controller *MoverController) *fyne.Container {
	cw := controller.width
	ch := controller.height
	rle := &RLE{}
	err := rle.Load("testdata/rats.rle")
	if err != nil {
		panic(err.Error())
	}
	cont := container.New(NewStaticLayout(cw, ch))
	coords := rle.coords
	lifeGen.AddCells(coords, lifeGen.CurrentGenId())
	timeText := NewMoverText("Time:", 10, 10, 10, fyne.TextAlignLeading)

	controller.SetOnUpdate(func(f float64) bool {
		timeText.SetText(fmt.Sprintf("Time: %05d Gen: %05d Cells:%05d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount))
		go lifeGen.NextGen()
		LifeResetDot()
		cell := lifeGen.CurrentGenRoot()
		gen := lifeGen.currentGenId
		for cell != nil {
			LifeGetDot(float32(cell.x), float32(cell.y), xOffset+(float32(gen)*100), yOffset+(float32(gen)*100), lifeGen.currentGenId, cont)
			cell = cell.next
		}
		LifeHideDot()
		return false
	})
	controller.AddMover(timeText, cont)
	return cont
}

func LifeHideDot() {
	for i := dotsPos; i < len(dots); i++ {
		dots[dotsPos].Hide()
	}
}

func LifeResetDot() {
	dotsPos = 0
}

func LifeGetDot(x, y, xOfs, yOfs float32, gen LifeGenId, container *fyne.Container) *canvas.Circle {
	if dotsPos >= len(dots) {
		for i := 0; i < 20; i++ {
			d := canvas.NewCircle(color.RGBA{0, 0, 255, 255})
			d.Hide()
			dots = append(dots, d)
			container.Add(d)
		}
	}
	dot := dots[dotsPos]
	dotsPos++
	dot.Position1 = fyne.Position{X: xOfs + (x * gridSize), Y: yOfs + (y * gridSize)}
	dot.Position2 = fyne.Position{X: xOfs + ((x * gridSize) + gridSize), Y: yOfs + ((y * gridSize) + gridSize)}
	dot.FillColor = genColor[gen]
	dot.Show()
	return dot
}
