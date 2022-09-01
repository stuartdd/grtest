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
	GetSpeed() (float32, float32)
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

/*
-------------------------------------------------------------------- MoverImage
*/
func NewMoverImage(x, y, w, h float32, image *canvas.Image) *MoverImage {
	return &MoverImage{imageSize: fyne.Size{Width: w, Height: h}, image: image, posx: x, posy: y, speedx: 0, speedy: 0}
}

func (mv *MoverImage) Update(time float32) {
	px := mv.posx + (mv.speedx * time)
	py := mv.posy + (mv.speedy * time)
	if (px != mv.posx) || (py != mv.posy) {
		mv.image.Move(fyne.Position{X: px, Y: py})
		mv.posx = px
		mv.posy = py
	}
}

func (mv *MoverImage) getCanvasObject() fyne.CanvasObject {
	return mv.image
}

func (mv *MoverImage) Init() {
	mv.image.Resize(mv.imageSize)
	mv.image.FillMode = canvas.ImageFillOriginal
}

func (mv *MoverImage) AdjustSpeed(x, y float32) {
	mv.speedx = mv.speedx + x
	mv.speedy = mv.speedy + y
}

func (mv *MoverImage) GetSpeed() (float32, float32) {
	return mv.speedx, mv.speedy
}

/*
-------------------------------------------------------------------- ControllerLayout
*/
func NewControllerContainer(width, height float32) *ControllerLayout {
	return &ControllerLayout{size: fyne.Size{Width: width, Height: height}, movers: make([]Movable, 0)}
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

/*
-------------------------------------------------------------------- main
*/
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

	controller.Add(moverImage1)
	container.Add(moverImage1.getCanvasObject())

	controller.Add(moverImage2)
	container.Add(moverImage2.getCanvasObject())

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
