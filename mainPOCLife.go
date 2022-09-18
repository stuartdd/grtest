package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
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

type MouseContainer struct {
	*fyne.Container
}

func NewMouseContainer(cw, ch float64) *fyne.Container {
	m := MouseContainer{&fyne.Container{Layout: &StaticLayout{size: fyne.Size{Width: float32(cw), Height: float32(ch)}}}}
	return m.Container
}

func (mc *MouseContainer) MouseDown(*desktop.MouseEvent) {
	fmt.Println("MouseDown")
}

func (mc *MouseContainer) MouseUp(*desktop.MouseEvent) {
	fmt.Println("MouseUp")
}

func (mc *MouseContainer) Tapped(*fyne.PointEvent) {
	fmt.Println("MouseUp")
}

var _ desktop.Mouseable = (*MouseContainer)(nil)
var _ fyne.Tappable = (*MouseContainer)(nil)

func POCLifeKeyPress(key *fyne.KeyEvent) {
	fmt.Println(key.Name)
	switch key.Name {
	case "Up":
		yOffset = yOffset + 50
	case "Down":
		yOffset = yOffset - 50
	case "Left":
		xOffset = xOffset + 50
	case "Right":
		xOffset = xOffset - 50
	case "=":
		gridSize = gridSize + 1
	case "-":
		if gridSize > 1 {
			gridSize = gridSize - 1
		}
	}
}

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, controller *MoverController) *fyne.Container {
	cw := controller.width
	ch := controller.height
	rle := &RLE{}
	err := rle.Load("testdata/Synth.rle")
	if err != nil {
		panic(err.Error())
	}
	cont := NewMouseContainer(cw, ch)

	//	cont := container.New(NewStaticLayout(cw, ch))
	coords := rle.coords
	timeText := NewMoverText("Time:", 10, 10, 20, fyne.TextAlignLeading)

	lifeGen = NewLifeGen(nil)

	lifeGen.AddCells(coords, lifeGen.CurrentGenId())

	controller.SetOnKeyPress(POCLifeKeyPress)
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
	dot.Resize(fyne.Size{Width: gridSize, Height: gridSize})
	dot.Show()

}