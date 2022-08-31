package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var runMainLoop = true

type Movable interface {
	// Move()
	Init()
	Update(time float32)
}

type MoverImage struct {
	speedx float32
	speedy float32
	posx   float32
	posy   float32

	iconSize fyne.Size
	icon     *canvas.Image
}

type ControllerLayout struct {
	size   fyne.Size
	movers []Movable
}

func (l *MoverImage) Update(time float32) {
	px := l.posx + (l.speedx * time)
	py := l.posy + (l.speedy * time)
	if (px != l.posx) || (py != l.posy) {
		l.icon.Move(fyne.Position{X: px, Y: py})
		l.posx = px
		l.posy = py
	}
}

func (l *MoverImage) Init() {
	l.icon.Resize(l.iconSize)
}

func (l *ControllerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.size = size
}

func (l *ControllerLayout) Update(time float32) {
	for _, m := range l.movers {
		m.Update(time)
	}
}

func (l *ControllerLayout) Init() {
	for _, m := range l.movers {
		m.Init()
	}
}

func (l *ControllerLayout) Add(m Movable) {
	l.movers = append(l.movers, m)
}

func (l *ControllerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return l.size
}

func main() {
	image := canvas.NewImageFromResource(Lander_Png)
	image.FillMode = canvas.ImageFillOriginal

	layout := &ControllerLayout{size: fyne.Size{Width: 500, Height: 300}, movers: make([]Movable, 0)}
	layout.Add(&MoverImage{iconSize: fyne.Size{Width: 40, Height: 40}, icon: image, posx: 0, posy: 0, speedx: 10, speedy: 10})

	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.SetCloseIntercept(func() {
		runMainLoop = false
		time.Sleep(time.Millisecond * 500)
		mainWindow.Close()
	})
	mainWindow.SetMaster()
	mainWindow.SetIcon(GoLogo_Png)

	cvs := container.New(layout, image)
	mainWindow.SetContent(cvs)
	go func() {
		layout.Init()
		timeOffset := time.Now().UnixMilli()
		for runMainLoop {
			time.Sleep(time.Millisecond * 100)
			layout.Update(float32(time.Now().UnixMilli()-timeOffset) / 1000.0)
			cvs.Refresh()
			timeOffset = time.Now().UnixMilli()
		}
		fmt.Println("CanRun No")
	}()

	mainWindow.ShowAndRun()
}
