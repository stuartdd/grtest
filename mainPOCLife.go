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
	gridSize    int64            = 6
	xOffset     int64            = 10
	yOffset     int64            = 10
	stopButton  *widget.Button
	startButton *widget.Button
	stepButton  *widget.Button
	timeText    = widget.NewLabel("")
	moverWidget *MoverWidget
	targetDot   *canvas.Circle
	rleFile     *RLE
	rleError    error

	FC_EMPTY = color.RGBA{255, 0, 0, 255}
	FC_ADDED = color.RGBA{0, 255, 0, 255}
	FC_FULL  = color.RGBA{0, 0, 255, 255}
	FC_CELL  = color.RGBA{0, 255, 255, 255}
)

func POCLifeMouseEvent(x, y float32, et MoverMouseEventType) {
	cellX, cellY := lifeScreenToCell(x, y)
	switch et {
	case MM_ME_TAP:
		c := lifeGen.GetCellFast(cellX, cellY)
		if c == 0 {
			lifeGen.AddCell(cellX, cellY, lifeGen.currentGenId)
			targetDot.FillColor = FC_ADDED
		} else {
			lifeGen.RemoveCell(cellX, cellY, lifeGen.currentGenId)
			targetDot.FillColor = FC_EMPTY
		}
		targetDot.Show()
	case MM_ME_MOVE:
		posX, posY := lifeCellToScreen(cellX, cellY)
		targetDot.Position1 = fyne.Position{X: posX, Y: posY}
		targetDot.Position2 = fyne.Position{X: posX + float32(gridSize), Y: posY + float32(gridSize)}
		targetDot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
		c := lifeGen.GetCellFast(cellX, cellY)
		if c == 0 {
			targetDot.FillColor = FC_EMPTY
		} else {
			targetDot.FillColor = FC_FULL
		}
		targetDot.Show()

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
	targetDot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
}

func POCLifeRunFor(n int) {
	lifeGen.SetRunFor(n, func(lg *LifeGen) {
		POCLifeStop()
	})
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_NONE)
	if mainController.animation != nil {
		mainController.animation.delay = 100
	}
	targetDot.Hide()
	stepButton.Disable()
	startButton.Disable()
	stopButton.Enable()
}

func POCLifeStop() {
	lifeGen.SetRunFor(0, nil)
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_MOVE|MM_ME_DOWN|MM_ME_UP|MM_ME_TAP)
	if mainController.animation != nil {
		mainController.animation.delay = 200
	}
	targetDot.Show()
	stepButton.Enable()
	startButton.Enable()
	stopButton.Disable()
}

/*
-------------------------------------------------------------------- main
*/
func MainPOCLife(mainWindow fyne.Window, width, height float64, controller *MoverController) *fyne.Container {
	controller.SetAnimationDelay(100)

	moverWidget = NewMoverWidget(width, height)
	targetDot = canvas.NewCircle(color.RGBA{250, 0, 0, 255})
	fmWidget := NewFileBrowserWidget(width, height, ".", "*.rle")
	fmWidget.Hide()

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

	topC.Add(lifeSeperator())
	topC.Add(widget.NewButton("File", func() {
		if fmWidget.Visible() {
			fmWidget.Hide()
		} else {
			POCLifeStop()
			fmWidget.SetOnMouseEvent(func(x, y float32, fbmet FileBrowseMouseEventType) {
				l := fmWidget.SelectByMouse(x, y)
				if l >= 0 {
					p := fmWidget.GetSelected()
					switch fbmet {
					case FB_ME_TAP:
						fmWidget.Refresh()
					case FB_ME_DTAP:
						if p == ".." {
							fmWidget.SetParentPath()
						} else {
							fmWidget.Hide()
							fmt.Printf("Selected %s", fmWidget.GetSelected())
							POCLifeStop()
							rleFile, rleError = NewRleFile(fmWidget.GetSelected())
							if rleError != nil {
								panic(rleError)
							}
							lifeGen.Clear()
							lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
							POCLifeRunFor(RUN_FOR_EVER)
							mainWindow.SetTitle(fmWidget.GetSelected())
						}
					}
				}
			}, FB_ME_TAP|FB_ME_DTAP)
			fmWidget.Show()
		}
	}))
	topC.Add(widget.NewButton("Restart", func() {
		POCLifeStop()
		lifeGen.Clear()
		lifeGen.AddCellsAtOffset(xOffset, yOffset, rleFile.coords, lifeGen.currentGenId)
	}))
	topC.Add(lifeSeperator())
	topC.Add(startButton)
	topC.Add(stopButton)
	topC.Add(stepButton)
	topC.Add(lifeSeperator())
	topC.Add(widget.NewButton("-", func() {
		POCLifeKeyPress("-")
	}))
	topC.Add(widget.NewButton("+", func() {
		POCLifeKeyPress("+")
	}))
	topC.Add(lifeSeperator())
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
	rleFile, rleError = NewRleFile("testdata/Infinite_growth.rle")
	if rleError != nil {
		panic(rleError)
	}
	lifeGen = NewLifeGen(nil, 0)
	lifeGen.AddCellsAtOffset(10, 10, rleFile.coords, lifeGen.currentGenId)
	POCLifeRunFor(RUN_FOR_EVER)
	mainWindow.SetTitle(fmt.Sprintf("File:%s", rleFile.fileName))

	controller.SetOnKeyPress(func(key *fyne.KeyEvent) {
		POCLifeKeyPress(string(key.Name))
	})

	controller.AddBeforeUpdate(func(f float64) bool {
		lifeGen.NextGen()
		POCLifeResetDot()
		gen := lifeGen.currentGenId

		cell := lifeGen.generations[lifeGen.currentGenId]
		for cell != nil {
			POCLifeGetDot(cell.x, cell.y, gen, moverWidget)
			cell = cell.next
		}
		timeText.SetText(fmt.Sprintf("Time: %05dms Gen: %05d Cells:%05d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount[lifeGen.currentGenId]))
		return false
	})
	moverWidget.AddTop(targetDot)
	moverWidget.SetFileBrowserWidget(fmWidget)

	return container.NewBorder(topC, botC, nil, nil, moverWidget)
}

func POCLifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

func POCLifeGetDot(x, y int64, gen LifeGenId, moverWidget *MoverWidget) {
	if dotsPos >= len(dots) {
		for i := 0; i < 20; i++ {
			d := canvas.NewCircle(color.RGBA{0, 0, 255, 255})
			d.Hide()
			dots = append(dots, d)
			moverWidget.AddBottom(d)
		}
	}
	dot := dots[dotsPos]
	dotsPos++
	//
	//	posX, posY := cellToScreen(x, y)
	//
	posX := float32((xOffset + x) * gridSize)
	posY := float32((yOffset + y) * gridSize)
	dot.Position1 = fyne.Position{X: posX, Y: posY}
	dot.Position2 = fyne.Position{X: posX + float32(gridSize), Y: posY + float32(gridSize)}
	dot.FillColor = FC_CELL
	dot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
	dot.Show()
}

func lifeCellToScreen(cellX, cellY int64) (float32, float32) {
	x := ((xOffset + cellX) * gridSize)
	y := ((yOffset + cellY) * gridSize)
	return float32(x), float32(y)
}

func lifeScreenToCell(mouseX, mouseY float32) (int64, int64) {
	cellX := int64((mouseX / float32(gridSize))) - xOffset
	cellY := int64((mouseY / float32(gridSize))) - yOffset
	return cellX, cellY
}

func lifeSeperator() *widget.Separator {
	sep := widget.NewSeparator()
	sep.Resize(fyne.Size{Width: 10, Height: sep.MinSize().Height})
	return sep
}
