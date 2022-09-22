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
	dots        []*canvas.Circle = make([]*canvas.Circle, 0)
	dotsPos     int              = 0
	gridSize    float32          = 4
	lifeGen     *LifeGen
	xOffset     float32       = 10
	yOffset     float32       = 10
	genColor    []color.Color = []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
	stopButton  *widget.Button
	startButton *widget.Button
	stepButton  *widget.Button
)

func POCLifeKeyPress(key string) {
	switch key {
	case "F1":
		if mainController.IsAnimation() {
			POCLifeStop()
		} else {
			POCLifeStart()
		}
	case "F2":
		if !mainController.IsAnimation() {
			mainController.Update(0)
		}
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

func POCLifeLoad(file string, err error) (*LifeGen, error) {
	if err != nil {
		return nil, err
	}
	rle := &RLE{}
	err = rle.Load(file)
	if err != nil {
		return nil, err
	}
	coords := rle.coords
	lg := NewLifeGen(nil)
	lg.AddCellsAtOffset(0, 0, coords, lg.currentGenId)
	minx, miny, maxx, maxy := lg.GetBounds()
	xOffset = float32(minx) + 10
	yOffset = float32(miny) + 10
	fmt.Printf("File: %s MinX: %4d MinY: %4d MaxX: %04d MaxY: %4d\n", file, minx, miny, maxx, maxy)
	return lg, nil
}

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, width, height float64, controller *MoverController) *fyne.Container {
	var err error
	moverWidget := NewMoverWidget(width, height)
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
	topC.Add(widget.NewSeparator())
	topC.Add(widget.NewButton("File", func() {
		go runMyFileDialog(mainWindow, "", func(file string, err error) {
			if err == nil {
				POCLifeStop()
				lifeGen, err = POCLifeLoad(file, nil)
				if err != nil {
					panic(err)
				}
				POCLifeStart()
			}
		})
	}))
	topC.Add(widget.NewSeparator())
	topC.Add(startButton)
	topC.Add(stopButton)
	topC.Add(stepButton)
	topC.Add(widget.NewSeparator())
	topC.Add(widget.NewButton("-", func() {
		POCLifeKeyPress("-")
	}))
	topC.Add(widget.NewButton("+", func() {
		POCLifeKeyPress("+")
	}))
	topC.Add(widget.NewSeparator())
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

	timeText := widget.NewLabel("")
	botC.Add(timeText)
	lifeGen, err = POCLifeLoad("testdata/1234_synth.rle", nil)
	if err != nil {
		panic(err)
	}

	controller.SetOnKeyPress(func(key *fyne.KeyEvent) {
		POCLifeKeyPress(string(key.Name))
	})

	controller.AddOnUpdateBefore(func(f float64) bool {
		lifeGen.NextGen()
		LifeResetDot()
		gen := lifeGen.currentGenId

		cell := lifeGen.generations[lifeGen.currentGenId]
		for cell != nil {
			LifeGetDot(float32(cell.x), float32(cell.y), xOffset, yOffset, gen, moverWidget)
			cell = cell.next
		}
		timeText.SetText(fmt.Sprintf("Time: %05dms Gen: %05d Cells:%05d", lifeGen.timeMillis, lifeGen.countGen, lifeGen.cellCount[lifeGen.currentGenId]))
		return true
	})

	return container.NewBorder(topC, botC, nil, nil, moverWidget)
}

func LifeResetDot() {
	dotsPos = 0
	for i := 0; i < len(dots); i++ {
		dots[i].Hide()
	}
}

func LifeGetDot(x, y, xOfs, yOfs float32, gen LifeGenId, container *MoverWidget) {
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
