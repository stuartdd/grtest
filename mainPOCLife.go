package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var (
	dots      []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos   int              = 0
	gridSize  float32          = 4
	lifeGen   *LifeGen
	xOffset   float32       = 10
	yOffset   float32       = 10
	genColor  []color.Color = []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
	countDots int           = 0
)

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, controller *MoverController) *fyne.Container {
	cw := controller.width
	ch := controller.height
	rle := &RLE{}
	err := rle.Load("testdata/reactions.rle")
	if err != nil {
		panic(err.Error())
	}
	cont := container.New(NewStaticLayout(cw, ch))
	coords := rle.coords
	timeText := NewMoverText("Time:", 10, 10, 20, fyne.TextAlignLeading)

	lifeGen = NewLifeGen(nil)

	lifeGen.AddCells(coords, lifeGen.CurrentGenId())
	controller.AddOnUpdate(func(f float64) bool {
		lifeGen.NextGen()
		LifeResetDot()
		countDots = 0
		gen := lifeGen.currentGenId

		cell := lifeGen.CurrentGenRoot()
		for cell != nil {
			LifeGetDot(float32(cell.x), float32(cell.y), xOffset, yOffset, gen, cont)
			cell = cell.next
			countDots++
		}
		timeText.SetText(fmt.Sprintf("Time: %05d Gen: %05d Cells:%05d DOTS:%d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount[lifeGen.currentGenId], countDots))
		return false
	})
	controller.AddMover(timeText, cont)
	return cont
}

func LifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

func LifeGetDot(x, y, xOfs, yOfs float32, gen LifeGenId, container *fyne.Container) {
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
	dot.FillColor = genColor[1]
	dot.Show()

}
