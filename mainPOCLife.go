package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var (
	lifeGen     *LifeGen
	dots        []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos     int              = 0
	gridSize    int64            = 6
	xOffset     int64            = 10
	yOffset     int64            = 10
	currentWd   string
	stopButton  *widget.Button
	startButton *widget.Button
	stepButton  *widget.Button
	timeText    = widget.NewLabel("")
	moverWidget *MoverWidget
	fbWidget    *FileBrowserWidget
	lifeWindow  fyne.Window
	targetDot   *canvas.Circle
	targetRect  *canvas.Rectangle
	rleFile     *RLE
	rleError    error

	FC_EMPTY  = color.RGBA{255, 0, 0, 255}   // Cell selector over an empty cell
	FC_ADDED  = color.RGBA{0, 255, 0, 255}   // Cell just added
	FC_FULL   = color.RGBA{0, 0, 255, 255}   // Cell selector over a cell
	FC_SELECT = color.RGBA{255, 255, 0, 255} // Cell colour inside selection rectangle
	FC_CELL   = color.RGBA{0, 255, 255, 255} // Normal, running cell colour

	COLOURS = []color.Color{FC_CELL, FC_SELECT, FC_FULL, FC_EMPTY} // Cell colour indexed by first two bits og the cell mode value
)

func POCLifeMouseEvent(me *MoverWidgetMouseEvent) {
	cellX1, cellY1 := lifeScreenToCell(float32(me.X1), float32(me.Y1))
	// cellX2, cellY2 := lifeScreenToCell(float32(me.X2), float32(me.Y2))
	switch me.Event {
	case MM_ME_TAP:
		c := lifeGen.GetCellFast(cellX1, cellY1)
		if c == 0 {
			lifeGen.AddCell(cellX1, cellY1, 0, lifeGen.currentGenId)
			targetDot.FillColor = FC_ADDED
		} else {
			lifeGen.RemoveCell(cellX1, cellY1, lifeGen.currentGenId)
			targetDot.FillColor = FC_EMPTY
		}
		targetDot.Show()
	case MM_ME_DTAP:
		if lifeGen.IsRunning() {
			POCLifeStop()
			return
		}
		POCLifeFile(cellX1, cellY1, false)
	case MM_ME_DRAG:
		targetRect.Move(*me.Position())
		targetRect.Resize(*me.Size())
		targetRect.Show()
		if me.Button == int(desktop.MouseButtonPrimary) {
			lifeGen.ClearMode(0b0)
		}
		x1, y1 := lifeScreenToCell(float32(me.X1), float32(me.Y1))
		x2, y2 := lifeScreenToCell(float32(me.X2), float32(me.Y2))
		lifeGen.CellsInBounds(x1, y1, x2, y2, func(lc *LifeCell) {
			lc.mode = lc.mode | 0x01
		})
	case MM_ME_MOVE:
		if me.Dragging {
			targetRect.Move(*me.Position())
			targetRect.Resize(*me.Size())
			targetRect.Show()
		} else {
			posX, posY := lifeCellToScreen(cellX1, cellY1)
			targetDot.Position1 = fyne.Position{X: posX, Y: posY}
			targetDot.Position2 = fyne.Position{X: posX + float32(gridSize), Y: posY + float32(gridSize)}
			targetDot.Resize(fyne.Size{Width: float32(gridSize), Height: float32(gridSize)})
			c := lifeGen.GetCellFast(cellX1, cellY1)
			if c == 0 {
				targetDot.FillColor = FC_EMPTY
			} else {
				targetDot.FillColor = FC_FULL
			}
			targetDot.Show()
			targetRect.Hide()
		}
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
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_DTAP)
	if mainController.animation != nil {
		mainController.animation.delay = 100
	}
	targetDot.Hide()
	targetRect.Hide()
	stepButton.Disable()
	startButton.Disable()
	stopButton.Enable()
}

func POCLifeFileZero() {
	POCLifeFile(xOffset, yOffset, true)
}

func POCLifeFile(cellPosX, cellPosY int64, clearCells bool) {
	if fbWidget.Visible() {
		fbWidget.Hide()
	} else {
		POCLifeStop()
		fbWidget.SetPath(currentWd)
		fbWidget.SetOnMouseEvent(func(x, y float32, fbmet FileBrowseMouseEventType) {
			l := fbWidget.SelectByMouse(x, y)
			if l >= 0 {
				p, selType := fbWidget.GetSelected()
				switch fbmet {
				case FB_ME_TAP:
					fbWidget.Refresh()
				case FB_ME_DTAP:
					switch selType {
					case FB_PARENT:
						fbWidget.SetParentPath()
					case FB_DIR:
						fbWidget.SetPath(p)
						currentWd = p
					case FB_FILE:
						fbWidget.Hide()
						POCLifeStop()
						rleFile, rleError = NewRleFile(p)
						if rleError != nil {
							panic(rleError)
						}
						if clearCells {
							lifeGen.Clear()
						}
						ofsx, ofsy := rleFile.RleCenter()
						lifeGen.AddCellsAtOffset(cellPosX-ofsx, cellPosY-ofsy, 0, rleFile.coords, lifeGen.currentGenId)
						POCLifeRunFor(RUN_FOR_EVER)
						lifeWindow.SetTitle(p)
					}
				}
			}
		}, FB_ME_TAP|FB_ME_DTAP)
		fbWidget.Show()
	}
}

func POCLifeStop() {
	lifeGen.SetRunFor(0, nil)
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_MOVE|MM_ME_DOWN|MM_ME_UP|MM_ME_TAP|MM_ME_DTAP)
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
	lifeWindow = mainWindow
	controller.SetAnimationDelay(100)
	currentWd, _ = os.Getwd()
	moverWidget = NewMoverWidget(width, height)
	targetDot = canvas.NewCircle(color.RGBA{250, 0, 0, 255})
	targetRect = &canvas.Rectangle{StrokeColor: color.RGBA{250, 0, 0, 255}, StrokeWidth: 1}
	fbWidget = NewFileBrowserWidget(width, height)
	fbWidget.Hide()
	fbWidget.SetOnFileFoundEvent(func(de fs.DirEntry, rootPath string, typ FileBrowserLineType) string {
		name := de.Name()
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			return ""
		}
		if typ == FB_DIR {
			return de.Name()
		}
		if strings.HasSuffix(strings.ToLower(name), ".rle") {
			rle, e := NewRleFile(path.Join(rootPath, name))
			if e != nil {
				return fmt.Sprintf("%s | %s", name, e.Error())
			} else {
				return fmt.Sprintf("%s | %s", name, rle.comment)
			}
		}
		return ""
	})
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
	topC.Add(widget.NewButton("File", POCLifeFileZero))
	topC.Add(widget.NewButton("Restart", func() {
		POCLifeStop()
		lifeGen.Clear()
		lifeGen.AddCellsAtOffset(xOffset, yOffset, 0, rleFile.coords, lifeGen.currentGenId)
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
	lifeGen.AddCellsAtOffset(10, 10, 0, rleFile.coords, lifeGen.currentGenId)
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
			POCLifeGetDot(cell.x, cell.y, cell.mode, gen, moverWidget)
			cell = cell.next
		}
		timeText.SetText(fmt.Sprintf("Time: %05dms Gen: %05d Cells:%05d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount[lifeGen.currentGenId]))
		return false
	})
	moverWidget.AddTop(targetDot)
	moverWidget.AddTop(targetRect)
	moverWidget.SetFileBrowserWidget(fbWidget)

	return container.NewBorder(topC, botC, nil, nil, moverWidget)
}

func POCLifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

func POCLifeGetDot(x, y int64, mode int, gen LifeGenId, moverWidget *MoverWidget) {
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
	dot.FillColor = COLOURS[mode&0x03]
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
