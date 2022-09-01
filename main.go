package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Movable interface {
	Init()
	Update(float32)
	AdjustSpeed(float32, float32)
	getCanvasObject() fyne.CanvasObject
}

type MoverImage struct {
	speedx float32
	speedy float32
	posx   float32
	posy   float32

	imageSize fyne.Size
	image     *canvas.Image
}

type ControllerLayout struct {
	size   fyne.Size
	movers []Movable
}

func NewControllerContainer(width, height float32) *ControllerLayout {
	return &ControllerLayout{size: fyne.Size{Width: width, Height: height}, movers: make([]Movable, 0)}
}

func NewMoverImage(x, y, w, h float32, image *canvas.Image) *MoverImage {
	return &MoverImage{imageSize: fyne.Size{Width: w, Height: h}, image: image, posx: x, posy: y, speedx: 0, speedy: 0}
}

func (l *MoverImage) Update(time float32) {
	px := l.posx + (l.speedx * time)
	py := l.posy + (l.speedy * time)
	if (px != l.posx) || (py != l.posy) {
		l.image.Move(fyne.Position{X: px, Y: py})
		l.posx = px
		l.posy = py
	}
}

func (l *MoverImage) getCanvasObject() fyne.CanvasObject {
	return l.image
}

func (l *MoverImage) Init() {
	l.image.Resize(l.imageSize)
	l.image.FillMode = canvas.ImageFillOriginal
}

func (l *MoverImage) AdjustSpeed(x, y float32) {
	l.speedx = l.speedx + x
	l.speedy = l.speedy + y
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

func (l *ControllerLayout) Add(m Movable, container *fyne.Container) {
	container.Add(m.getCanvasObject())
	l.movers = append(l.movers, m)
}

func (l *ControllerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return l.size
}

func main() {
	moverImage1 := NewMoverImage(0, 0, 40, 40, canvas.NewImageFromResource(Lander_Png))
	moverImage1.AdjustSpeed(10, 10)
	moverImage2 := NewMoverImage(0, 0, 20, 20, canvas.NewImageFromResource(Lander_Png))
	moverImage2.AdjustSpeed(5, 5)

	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Close()
	})
	mainWindow.SetMaster()
	mainWindow.SetIcon(GoLogo_Png)

	controller := NewControllerContainer(500, 500)
	container := container.New(controller)
	controller.Add(moverImage1, container)
	controller.Add(moverImage2, container)

	mainWindow.SetContent(container)

	controller.Init()
	var ft float32 = 0
	an := fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		if f == 1.0 {
			ft = 0
		} else {
			controller.Update(f - ft)
			container.Refresh()
			ft = f
		}
	}}
	an.Start()
	mainWindow.ShowAndRun()
	an.Stop()
}
