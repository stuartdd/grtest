package main

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const RUN_FOR_EVER = math.MaxInt

var (
	lifeGen     *LifeGen
	dots        []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos     int              = 0
	gridSize    int64            = 22
	xOffset     int64            = 10
	yOffset     int64            = 10
	selectX     int64
	selectY     int64
	genColor    []color.Color = []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
	stopButton  *widget.Button
	startButton *widget.Button
	stepButton  *widget.Button
	timeText    = widget.NewLabel("")
	moverWidget *MoverWidget
	lines       *MoverLines
	rleFile     *RLE
)

func POCLifeMouseEvent(x, y float32, et MoverMouseEventType) {
	gridX := int64(x / float32(gridSize))
	gridY := int64(y / float32(gridSize))
	gsW := lines.GetSizeAndCenter().Width / 2
	gsH := lines.GetSizeAndCenter().Height / 2
	switch et {
	case ME_TAP:
		fmt.Printf("TAP: %d, %d\n", selectX, selectY)
		lifeGen.AddCell(selectX, selectY, lifeGen.currentGenId)
	case ME_MOVE:
		selectX = gridX
		selectY = gridY
		fmt.Printf("MOVE: %d, %d\n", selectX, selectY)
		lines.SetCenter(float64(gridX*gridSize)+gsW, float64(gridY*gridSize)+gsH)
	}
}

func POCLifeKeyPress(key string) {
	switch key {
	case "F1":
		if lifeGen.IsRunning() {
			POCLifeStop()
		} else {
			POCLifeRunFor(RUN_FOR_EVER)
		}
		return
	case "F2":
		if !lifeGen.IsRunning() {
			POCLifeRunFor(1)
		}
		return
	case "Up":
		yOffset = yOffset + 50
	case "Down":
		yOffset = yOffset - 50
	case "Left":
		xOffset = xOffset + 50
	case "Right":
		xOffset = xOffset - 50
	case "=", "+":
		POCLifeSetGridSize(true)
	case "-", "_":
		POCLifeSetGridSize(false)
	}
}

func POCLifeSetGridSize(inc bool) {
	if inc {
		gridSize++
	} else {
		if gridSize > 1 {
			gridSize--
		}
	}
	lines.SetSize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
}

func POCLifeRunFor(n int) {
	lifeGen.SetRunFor(n, func(lg *LifeGen) {
		POCLifeStop()
	})
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, ME_NONE)
	lines.SetVisible(false)
	stepButton.Disable()
	startButton.Disable()
	stopButton.Enable()
}

func POCLifeStop() {
	lifeGen.SetRunFor(0, nil)
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, ME_MOVE|ME_DOWN|ME_UP|ME_TAP)
	lines.SetVisible(true)
	stepButton.Enable()
	startButton.Enable()
	stopButton.Disable()
}

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, width, height float64, controller *MoverController) *fyne.Container {
	moverWidget = NewMoverWidget(width, height)
	lines = NewMoverRect(color.RGBA{250, 0, 0, 255}, 200, 200, float64(gridSize), float64(gridSize), 0)

	topC := container.NewHBox()
	botC := container.NewPadded()
	startButton = widget.NewButton("Start (F1)", func() {
		POCLifeRunFor(RUN_FOR_EVER)
	})
	stepButton = widget.NewButton("Step (F2)", func() {
		POCLifeRunFor(1)
	})
	stopButton = widget.NewButton("Stop (F1)", func() {
		POCLifeStop()
	})
	stepButton.Disable()
	startButton.Disable()

	topC.Add(widget.NewButton("Close (Esc)", func() {
		mainWindow.Close()
	}))

	topC.Add(seperator())
	topC.Add(widget.NewButton("File", func() {
		go runMyFileDialog(mainWindow, "", func(file string, err error) {
			if err == nil {
				POCLifeStop()
				rleFile, err = NewRleFile(file)
				if err != nil {
					panic(err)
				}
				lifeGen.Clear()
				lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
			}
		})
	}))
	topC.Add(widget.NewButton("Restart", func() {
		POCLifeStop()
		lifeGen.Clear()
		lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
	}))
	topC.Add(seperator())
	topC.Add(startButton)
	topC.Add(stopButton)
	topC.Add(stepButton)
	topC.Add(seperator())
	topC.Add(widget.NewButton("-", func() {
		POCLifeKeyPress("-")
	}))
	topC.Add(widget.NewButton("+", func() {
		POCLifeKeyPress("+")
	}))
	topC.Add(seperator())
	topC.Add(widget.NewButton("<", func() {
		POCLifeKeyPress("Left")
	}))
	topC.Add(widget.NewButton("^", func() {
		POCLifeKeyPress("Up")
	}))
	topC.Add(widget.NewButton("v", func() {
		POCLifeKeyPress("Down")
	}))
	topC.Add(widget.NewButton(">", func() {
		POCLifeKeyPress("Right")
	}))

	botC.Add(timeText)

	var err error
	rleFile, err = NewRleFile("testdata/1234_synth.rle")
	if err != nil {
		panic(err)
	}
	lifeGen = NewLifeGen(nil)
	lifeGen.AddCellsAtOffset(10, 10, rleFile.coords, lifeGen.currentGenId)
	POCLifeRunFor(RUN_FOR_EVER)
	mainWindow.SetTitle(fmt.Sprintf("File:%s", rleFile.fileName))

	controller.SetOnKeyPress(func(key *fyne.KeyEvent) {
		POCLifeKeyPress(string(key.Name))
	})

	controller.AddOnUpdateBefore(func(f float64) bool {
		lifeGen.NextGen()
		onUpdateBefore()
		return false
	})
	moverWidget.AddMover(lines)
	return container.NewBorder(topC, botC, nil, nil, moverWidget)
}

func onUpdateBefore() {
	LifeResetDot()
	gen := lifeGen.currentGenId

	cell := lifeGen.generations[lifeGen.currentGenId]
	for cell != nil {
		LifeGetDot(cell.x, cell.y, xOffset, yOffset, gen, moverWidget)
		cell = cell.next
	}
	timeText.SetText(fmt.Sprintf("Time: %05dms Gen: %05d Cells:%05d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount[lifeGen.currentGenId]))
}

func LifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

// func dotTo(c *canvas.Circle, x, y float32, show bool) {
// 	gs2 := float32(gridSize / 2)
// 	c.Position1 = fyne.Position{X: x - gs2, Y: y - gs2}
// 	c.Position2 = fyne.Position{X: x + gs2, Y: y + gs2}
// }

func LifeGetDot(x, y, xOfs, yOfs int64, gen LifeGenId, moverWidget *MoverWidget) {
	if dotsPos >= len(dots) {
		for i := 0; i < 20; i++ {
			d := canvas.NewCircle(color.RGBA{0, 0, 255, 255})
			d.Hide()
			dots = append(dots, d)
			moverWidget.Add(d)
		}
	}
	gs2 := float32(gridSize / 2)
	x = xOfs + (x * int64(gridSize))
	y = xOfs + (y * int64(gridSize))
	dot := dots[dotsPos]
	dotsPos++
	dot.Position1 = fyne.Position{X: float32(x) - gs2, Y: float32(y) - gs2}
	dot.Position2 = fyne.Position{X: float32(x) + gs2, Y: float32(y) + gs2}
	dot.FillColor = genColor[1]
	dot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
	dot.Show()
}

func seperator() *widget.Separator {
	sep := widget.NewSeparator()
	sep.Resize(fyne.Size{Width: 10, Height: sep.MinSize().Height})
	return sep
}
