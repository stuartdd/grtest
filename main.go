package main

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Movable interface {
	Init()
	Update(float64)
	SetPosition(float64, float64)
	GetPosition() (float64, float64)
	SetSpeed(float64, float64)
	GetSpeed() (float64, float64)
	GetCanvasObject() fyne.CanvasObject
}

type MoverLines struct {
	speedx   float64
	accx     float64
	speedy   float64
	accy     float64
	speedRot float64
	accRot   float64
	centerX  float64
	centerY  float64
	lines    []*canvas.Line
}

type MoverImage struct {
	speedx float64
	speedy float64
	posx   float64
	posy   float64

	imageSize fyne.Size
	image     *canvas.Image
}

type ControllerLayout struct {
	size   fyne.Size
	movers []Movable
}

var _ Movable = (*MoverLines)(nil)
var _ Movable = (*MoverImage)(nil)

/*
-------------------------------------------------------------------- MoverLines
*/

func NewMoverLines(centerX, centerY, speedRot float64) *MoverLines {
	return &MoverLines{speedx: 0, speedy: 0, speedRot: speedRot, centerX: centerX, centerY: centerY, lines: make([]*canvas.Line, 0)}
}

func (mv *MoverLines) Update(time float64) {
	mv.accx = mv.accx + (mv.speedx * time)
	var dx int = 0
	if math.Abs(mv.accx) > 1 {
		dx = int(mv.accx)
		mv.accx = mv.accx - float64(dx)
	}
	mv.accy = mv.accy + (mv.speedy * time)
	var dy int = 0
	if math.Abs(mv.accy) > 1 {
		dy = int(mv.accy)
		mv.accy = mv.accy - float64(dy)
	}
	mv.accRot = mv.accRot + (mv.speedRot * float64(time))
	var ra int = 0
	if math.Abs(mv.accRot) > 1 {
		ra = int(mv.accRot)
		mv.accRot = mv.accRot - float64(ra)
	}

	for _, l := range mv.lines {
		if dx != 0 {
			l.Position1.X = l.Position1.X + float32(dx)
			l.Position2.X = l.Position2.X + float32(dx)
		}
		if dy != 0 {
			l.Position1.Y = l.Position1.Y + float32(dy)
			l.Position2.Y = l.Position2.Y + float32(dy)
		}
		if ra != 0 {
			rotatePoint(mv.centerX, mv.centerY, &l.Position1, ra)
			rotatePoint(mv.centerX, mv.centerY, &l.Position2, ra)
		}
	}
	mv.centerX = mv.centerX + float64(dx)
	mv.centerY = mv.centerY + float64(dy)
}

func (mv *MoverLines) GetCanvasObject() fyne.CanvasObject {
	container := container.New(&ControllerLayout{})
	for _, l := range mv.lines {
		container.Add(l)
	}
	return container
}

func (mv *MoverLines) Init() {
}

func (mv *MoverLines) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverLines) SetPosition(x, y float64) {
	dx := x - mv.centerX
	dy := y - mv.centerY
	for _, l := range mv.lines {
		l.Position1.X = l.Position1.X + float32(dx)
		l.Position2.X = l.Position2.X + float32(dx)
		l.Position1.Y = l.Position1.Y + float32(dy)
		l.Position2.Y = l.Position2.Y + float32(dy)
	}
	mv.centerX = x
	mv.centerY = y
}

func (mv *MoverLines) GetPosition() (float64, float64) {
	return mv.centerX, mv.centerY
}

func (mv *MoverLines) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

func (mv *MoverLines) AddLine(colour color.Color, x1, y1, x2, y2 float32) {
	line := canvas.NewLine(colour)
	line.Position1.X = x1
	line.Position1.Y = y1
	line.Position2.X = x2
	line.Position2.Y = y2
	mv.lines = append(mv.lines, line)
}
func (mv *MoverLines) AddLines(colour color.Color, xy ...float32) {
	for i := 0; i < len(xy); i = i + 2 {
		mv.AddLineToo(colour, xy[i], xy[i+1])
	}
}

func (mv *MoverLines) AddLineToo(colour color.Color, x2, y2 float32) {
	line := canvas.NewLine(colour)
	var x1 float32 = 0.0
	var y1 float32 = 0.0
	le := len(mv.lines)
	if le > 0 {
		lf := mv.lines[le-1]
		x1 = lf.Position2.X
		y1 = lf.Position2.Y
	}
	line.Position1.X = x1
	line.Position1.Y = y1
	line.Position2.X = x2
	line.Position2.Y = y2
	mv.lines = append(mv.lines, line)
}

/*
-------------------------------------------------------------------- MoverImage
*/
func NewMoverImage(x, y, w, h float64, image *canvas.Image) *MoverImage {
	return &MoverImage{imageSize: fyne.Size{Width: float32(w), Height: float32(h)}, image: image, posx: x, posy: y, speedx: 0, speedy: 0}
}

func (mv *MoverImage) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if (dx != 0) || (dy != 0) {
		mv.posx = mv.posx + dx
		mv.posy = mv.posy + dy
		mv.image.Move(fyne.Position{X: float32(mv.posx), Y: float32(mv.posy)})
	}
}

func (mv *MoverImage) GetCanvasObject() fyne.CanvasObject {
	return mv.image
}

func (mv *MoverImage) Init() {
	mv.image.Resize(mv.imageSize)
	mv.image.FillMode = canvas.ImageFillOriginal
}

func (mv *MoverImage) SetPosition(x, y float64) {
	mv.posx = x
	mv.posy = y
	mv.image.Move(fyne.Position{X: float32(mv.posx), Y: float32(mv.posy)})
}

func (mv *MoverImage) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverImage) GetPosition() (float64, float64) {
	return float64(mv.image.Position().X), float64(mv.image.Position().Y)
}

func (mv *MoverImage) GetSpeed() (float64, float64) {
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

func (l *ControllerLayout) Update(time float64) {
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
	// testCos()
	// testSine()
	// os.Exit(0)
	moverImage1 := NewMoverImage(0, 0, 40, 40, canvas.NewImageFromResource(Lander_Png))
	moverImage1.SetSpeed(100, 100)
	moverImage1.SetPosition(100, 0)
	moverImage2 := NewMoverImage(0, 0, 20, 20, canvas.NewImageFromResource(Lander_Png))
	moverImage2.SetSpeed(5, 5)

	lines1 := NewMoverLines(150, 150, -9)
	lines1.AddLine(color.White, 100, 100, 150, 150)
	lines1.AddLineToo(color.White, 200, 100)
	lines1.SetSpeed(5, 5)

	lines2 := NewMoverLines(50, 50, 20)
	lines2.AddLines(color.White, 0, 0, 100, 0, 100, 100, 0, 100, 0, 0)
	lines2.SetSpeed(5, 5)
	lines2.SetPosition(100, 50)

	lines3 := NewMoverLines(0, 0, 0)
	lines3.AddLineToo(color.White, 400, 400)

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
	container.Add(moverImage1.GetCanvasObject())

	controller.Add(moverImage2)
	container.Add(moverImage2.GetCanvasObject())

	controller.Add(lines1)
	container.Add(lines1.GetCanvasObject())
	controller.Add(lines2)
	container.Add(lines2.GetCanvasObject())
	controller.Add(lines3)
	container.Add(lines3.GetCanvasObject())

	mainWindow.SetContent(container)
	controller.Init()

	var ft float32 = 0
	an := fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		controller.Update(float64(f - ft))
		if f == 1.0 {
			ft = 0
		} else {
			ft = f
		}
		container.Refresh()
	}}
	an.Start()
	mainWindow.ShowAndRun()
	an.Stop()
}
