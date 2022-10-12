package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	SELECT_MODE_MASK = 0b00000001
	COLOUR_MODE_MASK = 0b00000011
)

var (
	moverWidget     *MoverWidget
	fbWidget        *FileBrowserWidget
	lifeWindow      fyne.Window
	lifeController  *MoverController
	lifeGen         *LifeGen
	lifeGenStopped  bool
	selectedCellsXY []int64

	dots         []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos      int              = 0
	gridSize     int64            = 6
	xOffset      int64            = 0
	yOffset      int64            = 0
	currentDelay int64            = 100
	currentWd    string
	stopButton   *widget.Button
	startButton  *widget.Button
	stepButton   *widget.Button
	deleteButton *widget.Button
	saveButton   *widget.Button
	clearButton  *widget.Button
	fasterButton *widget.Button
	slowerButton *widget.Button
	timeText     = widget.NewLabel("")
	targetDot    *canvas.Circle
	targetRect   *canvas.Rectangle
	rleFile      *RLE
	rleError     error

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
		c := lifeGen.GetCell(cellX1, cellY1)
		if me.Button == int(desktop.MouseButtonPrimary) {
			if c == 0 {
				lifeGen.AddCell(cellX1, cellY1, 0)
				targetDot.FillColor = FC_ADDED
			} else {
				lifeGen.RemoveCell(cellX1, cellY1)
				targetDot.FillColor = FC_EMPTY
			}
		} else {
			if len(selectedCellsXY) > 0 {
				lifeGen.AddCellsAtOffset(cellX1, cellY1, 0, selectedCellsXY)
			}
		}
		targetDot.Show()
	case MM_ME_DTAP:
		POCLifeStop()
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
			lc.mode = lc.mode | SELECT_MODE_MASK
		})
		tmp := lifeGen.ListCellsWithMode(SELECT_MODE_MASK)
		if len(tmp) > 0 {
			selectedCellsXY, _, _ = POCNormaliseCoords(tmp)
		}
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
			c := lifeGen.GetCell(cellX1, cellY1)
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
	case "s", "S":
		POCLifeSetSlower()
	case "f", "F":
		POCLifeSetFaster()
	case "c", "C":
		POCLifeHome()
	}
}
func POCLifeHome() {
	runsRemaining := POCLifeStop()
	midX := (int64(lifeWindow.Canvas().Size().Width) / gridSize)
	midY := (int64(lifeWindow.Canvas().Size().Height) / gridSize)
	x1, y1, x2, y2 := lifeGen.GetBounds()
	xOffset = ((midX - (x2 - x1)) / 2) - x1
	yOffset = ((midY - (y2 - y1)) / 2) - y1
	if runsRemaining > 0 {
		POCLifeRunFor(runsRemaining)
	}
}

func POCLifeSetSlower() {
	currentDelay = currentDelay + 10
	fasterButton.Enable()
	if currentDelay >= 400 {
		currentDelay = 400
		slowerButton.Disable()
	}
	lifeController.SetAnimationDelay(currentDelay)
}

func POCLifeSetFaster() {
	currentDelay = currentDelay - 10
	slowerButton.Enable()
	if currentDelay <= 10 {
		currentDelay = 10
		fasterButton.Disable()
	}
	lifeController.SetAnimationDelay(currentDelay)
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

func POCLifeStop() int {
	runsRemaining := lifeGen.GetRunFor()
	if lifeGen.IsRunning() {
		notStopped := true
		lifeGen.SetRunFor(1, func(lg *LifeGen) {
			notStopped = false
		})
		for notStopped {
			time.Sleep(time.Millisecond * 100)
		}
	}
	lifeGenStopped = true
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_MOVE|MM_ME_DOWN|MM_ME_UP|MM_ME_TAP|MM_ME_DTAP)

	lifeController.SetAnimationDelay(200)
	targetDot.Show()
	stepButton.Enable()
	clearButton.Enable()
	startButton.Enable()
	stopButton.Disable()
	return runsRemaining
}

func POCLifeRunFor(n int) {
	lifeGen.SetRunFor(n, func(lg *LifeGen) {
		POCLifeStop()
	})
	lifeGenStopped = false
	moverWidget.SetOnMouseEvent(POCLifeMouseEvent, MM_ME_DTAP)
	lifeController.SetAnimationDelay(currentDelay)
	targetDot.Hide()
	targetRect.Hide()
	stepButton.Disable()
	clearButton.Disable()
	startButton.Disable()
	stopButton.Enable()
	deleteButton.Hide()
	saveButton.Hide()
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
		fbWidget.SetOnSelectedEvent(func(fil, path string) error {
			POCLifeStop()
			rleFile, rleError = NewRleFile(fil)
			if rleError != nil {
				panic(rleError)
			}
			currentWd = path
			if clearCells {
				lifeGen.Reset()
			}
			ofsx, ofsy := rleFile.Center()
			lifeGen.AddCellsAtOffset(cellPosX-ofsx, cellPosY-ofsy, 0, rleFile.coords)
			POCLifeRunFor(RUN_FOR_EVER)
			lifeWindow.SetTitle(fil)
			return nil
		})
		fbWidget.Show()
	}
}

/*
-------------------------------------------------------------------- main
*/
func MainPOCLife(mainWindow fyne.Window, width, height float64, moverController *MoverController) *fyne.Container {
	lifeWindow = mainWindow
	lifeController = moverController
	lifeController.SetAnimationDelay(100)
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
	topV := container.NewVBox()
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
	clearButton = widget.NewButton("Clear", func() {
		lifeGen.Reset()
	})
	deleteButton = widget.NewButton("Delete", func() {
		if len(selectedCellsXY) > 0 {
			lifeGen.RemoveCellsWithMode(SELECT_MODE_MASK)
		}
	})
	saveButton = widget.NewButton("Save", func() {
		if len(selectedCellsXY) > 0 {
			fbWidget.SaveFormShow(func(s string, save bool, err error) error {
				fmt.Printf("SAVE %t %s\n", save, s)
				fbWidget.SaveFormHide()
				fbWidget.Hide()
				return err
			})
			POCLifeFileZero()
		}
	})
	fasterButton = widget.NewButton("F", func() {
		POCLifeSetFaster()
	})
	slowerButton = widget.NewButton("S", func() {
		POCLifeSetSlower()
	})
	deleteButton.Hide()
	saveButton.Hide()
	stepButton.Disable()
	startButton.Disable()
	clearButton.Disable()

	topC.Add(widget.NewButton("Close (Esc)", func() {
		mainWindow.Close()
	}))
	topC.Add(lifeSeperator())
	topC.Add(widget.NewButton("File", POCLifeFileZero))
	topC.Add(widget.NewButton("Restart", func() {
		POCLifeStop()
		lifeGen.Reset()
		lifeGen.AddCellsAtOffset(xOffset, yOffset, 0, rleFile.coords)
	}))
	topC.Add(lifeSeperator())
	topC.Add(startButton)
	topC.Add(stopButton)
	topC.Add(stepButton)
	topC.Add(clearButton)
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
	topC.Add(lifeSeperator())
	topC.Add(widget.NewButton("C", func() {
		POCLifeKeyPress("C")
	}))
	topC.Add(fasterButton)
	topC.Add(slowerButton)
	topC.Add(lifeSeperator())
	topC.Add(deleteButton)
	topC.Add(saveButton)

	botC.Add(timeText)
	rleFile, rleError = NewRleFile("testdata/Infinite_growth.rle")
	if rleError != nil {
		panic(rleError)
	}
	lifeGen = NewLifeGen(nil, 0)
	lifeGen.AddCellsAtOffset(10, 10, 0, rleFile.coords)
	POCLifeRunFor(RUN_FOR_EVER)
	mainWindow.SetTitle(fmt.Sprintf("File:%s", rleFile.fileName))

	lifeController.SetOnKeyPress(func(key *fyne.KeyEvent) {
		POCLifeKeyPress(string(key.Name))
	})

	lifeController.AddBeforeUpdate(func(f float64) bool {
		if lifeGenStopped {
			if len(selectedCellsXY) > 0 {
				if !saveButton.Visible() {
					saveButton.Show()
				}
			} else {
				if saveButton.Visible() {
					saveButton.Hide()
				}
			}
			if lifeGen.CountCellsWithMode(SELECT_MODE_MASK) > 0 {
				if !deleteButton.Visible() {
					deleteButton.Show()
				}
			} else {
				if deleteButton.Visible() {
					deleteButton.Hide()
				}
			}

		}
		lifeGen.NextGen()
		POCLifeResetDot()
		cell := lifeGen.GetRootCell()
		for cell != nil {
			POCLifeGetDot(cell.x, cell.y, cell.mode, moverWidget)
			cell = cell.next
		}
		timeText.SetText(fmt.Sprintf("Delay: %03dms Time: %05dms Gen: %05d Cells:%05d", lifeController.GetAnimationDelay(), lifeGen.GetGenerationTime(), lifeGen.GetGenerationCount(), lifeGen.GetCellCount()))
		return false
	})
	moverWidget.AddTop(targetDot)
	moverWidget.AddTop(targetRect)
	moverWidget.SetFileBrowserWidget(fbWidget)

	topV.Add(topC)
	topV.Add(fbWidget.InputForm("Save Selected Cells to a RLE File"))
	return container.NewBorder(topV, botC, nil, nil, moverWidget)
}

func POCLifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

func POCLifeGetDot(x, y int64, mode int, moverWidget *MoverWidget) {
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
	dot.FillColor = COLOURS[mode&COLOUR_MODE_MASK]
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
