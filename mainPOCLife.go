package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	lifeGen     *LifeGen
	dots        []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos     int              = 0
	gridSize    int64            = 4
	xOffset     int64            = 10
	yOffset     int64            = 10
	genColor    []color.Color    = []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
	stopButton  *widget.Button
	startButton *widget.Button
	stepButton  *widget.Button
	timeText    = widget.NewLabel("")
	moverWidget *MoverWidget
	rleFile     *RLE
)

func POCLifeKeyPress(key string) {
	switch key {
	case "F1":
		if mainController.IsAnimation() {
			POCLifeStop()
		} else {
			POCLifeStart()
		}
		return
	case "F2":
		if !mainController.IsAnimation() {
			mainController.Update(0)
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
		gridSize = gridSize + 1
	case "-", "_":
		if gridSize > 1 {
			gridSize = gridSize - 1
		}
	}
	onUpdateBefore()
}

func POCLifeStart() {
	mainController.StartAnimation()
	stepButton.Disable()
	startButton.Disable()
	stopButton.Enable()
}

func POCLifeStop() {
	mainController.StopAnimation()
	stepButton.Enable()
	startButton.Enable()
	stopButton.Disable()
}

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, width, height float64, controller *MoverController) *fyne.Container {
	moverWidget = NewMoverWidget(width, height)
	topC := container.NewHBox()
	botC := container.NewPadded()
	startButton = widget.NewButton("Start (F1)", func() {
		POCLifeStart()
	})
	stepButton = widget.NewButton("Step (F2)", func() {
		controller.Update(0)
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
				onUpdateBefore()
			}
		})
	}))
	topC.Add(widget.NewButton("Restart", func() {
		if mainController.IsAnimation() {
			POCLifeStop()
			lifeGen.Clear()
			lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
			POCLifeStart()
		} else {
			lifeGen.Clear()
			lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
			onUpdateBefore()
		}
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
	mainWindow.SetTitle(fmt.Sprintf("File:%s", rleFile.fileName))

	controller.SetOnKeyPress(func(key *fyne.KeyEvent) {
		POCLifeKeyPress(string(key.Name))
	})

	controller.AddOnUpdateBefore(func(f float64) bool {
		lifeGen.NextGen()
		onUpdateBefore()
		return true
	})

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

func LifeGetDot(x, y, xOfs, yOfs int64, gen LifeGenId, container *MoverWidget) {
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
	dot.Position1 = fyne.Position{X: float32(xOfs + (x * gridSize)), Y: float32(yOfs + (y * gridSize))}
	dot.Position2 = fyne.Position{X: float32(xOfs + ((x * gridSize) + gridSize)), Y: float32(yOfs + ((y * gridSize) + gridSize))}
	dot.FillColor = genColor[1]
	dot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
	dot.Show()
}

func seperator() *widget.Separator {
	sep := widget.NewSeparator()
	sep.Resize(fyne.Size{Width: 10, Height: sep.MinSize().Height})
	return sep
}
