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
	GetSize() *SizeAndCenter
	SetSize(fyne.Size)
}

type SizeAndCenter struct {
	Width   float32
	Height  float32
	CenterX float32
	CenterY float32
}

type MoverLines struct {
	speedx     float64
	accx       float64
	speedy     float64
	accy       float64
	currentAng int
	accAng     float64
	speedAng   float64
	centerX    float64
	centerY    float64
	lines      []*canvas.Line
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

func NewSizeAndCenter(w, h, x, y float32) *SizeAndCenter {
	return &SizeAndCenter{Width: w, Height: h, CenterX: x, CenterY: y}
}

func (sc *SizeAndCenter) Size() fyne.Size {
	return fyne.Size{Height: sc.Height, Width: sc.Width}
}

func (sc *SizeAndCenter) Center() fyne.Position {
	return fyne.Position{X: sc.CenterX, Y: sc.CenterY}
}

func NewMoverRect(colour color.Color, centerX, centerY, w, h, speedAng float64) *MoverLines {
	ml := &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerX: centerX, centerY: centerY, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
	ml.AddLine(colour, float32(centerX-w/2), float32(centerY-h/2), float32(centerX+w/2), float32(centerY-h/2))
	ml.AddLineToo(colour, float32(centerX+w/2), float32(centerY+h/2))
	ml.AddLineToo(colour, float32(centerX-w/2), float32(centerY+h/2))
	ml.AddLineToo(colour, float32(centerX-w/2), float32(centerY-h/2))
	return ml
}

/*
-------------------------------------------------------------------- MoverLines
*/

func NewMoverLines(centerX, centerY, speedAng float64) *MoverLines {
	return &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerX: centerX, centerY: centerY, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
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
	// Adjust the angle keeping it within bounds
	mv.accAng = mv.accAng + (mv.speedAng * float64(time))
	if mv.accAng > 360.0 {
		mv.accAng = mv.accAng - 360.0
	}
	if mv.accAng < 0 {
		mv.accAng = mv.accAng + 360.0
	}
	// Calc how much we need to rotate!
	intAng := int(mv.accAng)
	ra := 0
	if mv.currentAng != intAng {
		ra = intAng - mv.currentAng
		mv.currentAng = intAng
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

func (mv *MoverLines) GetAngle() int {
	return mv.currentAng
}

func (mv *MoverLines) SetSize(size fyne.Size) {
	currentSize := mv.GetSize()
	scaleX := float64(size.Width / currentSize.Width)
	scaleY := float64(size.Height / currentSize.Height)
	for _, l := range mv.lines {
		scalePoint(mv.centerX, mv.centerY, &l.Position1, scaleX, scaleY)
		scalePoint(mv.centerX, mv.centerY, &l.Position2, scaleX, scaleY)
	}
}

func (mv *MoverLines) GetSize() *SizeAndCenter {
	maxx := minFloat32
	maxy := minFloat32
	minx := maxFloat32
	miny := maxFloat32
	for _, l := range mv.lines {
		p1 := l.Position1
		p2 := l.Position2
		if p1.X > maxx {
			maxx = p1.X
		}
		if p1.X < minx {
			minx = p1.X
		}
		if p2.X > maxx {
			maxx = p2.X
		}
		if p2.X < minx {
			minx = p2.X
		}

		if p1.Y > maxy {
			maxy = p1.Y
		}
		if p1.Y < miny {
			miny = p1.Y
		}
		if p2.Y > maxy {
			maxy = p2.Y
		}
		if p2.Y < miny {
			miny = p2.Y
		}
	}
	w := maxx - minx
	h := maxy - miny
	return NewSizeAndCenter(w, h, (minx + w/2), (miny + h/2))
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

func (mv *MoverImage) GetSize() *SizeAndCenter {
	return NewSizeAndCenter(mv.imageSize.Width, mv.imageSize.Height, 0, 0)
}

func (mv *MoverImage) SetSize(size fyne.Size) {
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
	moverImage2 := NewMoverImage(0, 0, 40, 40, canvas.NewImageFromResource(Lander_Png))
	moverImage2.SetSpeed(5, 5)

	lines1 := NewMoverLines(200, 200, 10)
	lines1.AddLine(color.White, 200, 200, 200, 300)
	lines1.AddLineToo(color.White, 100, 300)
	lines1.AddLineToo(color.White, 200, 200)
	lines1.SetSpeed(2, 2)

	rect := NewMoverRect(color.RGBA{250, 0, 0, 255}, 200, 200, 100, 100, 0)

	lines2 := NewMoverLines(50, 50, 50)
	lines2.AddLines(color.RGBA{0, 255, 0, 150}, 0, 0, 100, 0, 100, 100, 0, 100, 0, 0)
	lines2.SetSpeed(5, 5)

	lines3 := NewMoverLines(0, 0, 0)
	lines3.AddLineToo(color.White, 1000, 1000)

	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Close()
	})
	mainWindow.SetMaster()
	mainWindow.SetIcon(GoLogo_Png)

	controller := NewControllerContainer(500, 500)
	container := container.New(controller)

	controller.Add(moverImage2)
	container.Add(moverImage2.GetCanvasObject())

	controller.Add(lines1)
	container.Add(lines1.GetCanvasObject())
	controller.Add(lines2)
	container.Add(lines2.GetCanvasObject())
	controller.Add(lines3)
	container.Add(lines3.GetCanvasObject())
	// controller.Add(rect)
	container.Add(rect.GetCanvasObject())

	mainWindow.SetContent(container)
	controller.Init()

	an := startAnimation(controller, container)
	go func() {
		//		sec := 0
		//		la := 0
		for {
			time.Sleep(time.Millisecond * 100)
			s := lines1.GetSize()
			rect.SetSize(s.Size())
			rect.SetPosition(float64(s.CenterX), float64(s.CenterY))
			//			a := lines1.GetAngle()
			//			aa := lines1.accAng
			//			fmt.Printf("Time %d: (aa=%f la=%d) a=%d w=%f h=%f\n", sec, aa, a-la, a, s.Width, s.Height)
			//			la = a
			//			sec++
		}
	}()
	mainWindow.ShowAndRun()
	an.Stop()
}

func startAnimation(controller *ControllerLayout, container *fyne.Container) *fyne.Animation {
	var ft float32 = 0
	an := &fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		controller.Update(float64(f - ft))
		if f == 1.0 {
			ft = 0
		} else {
			ft = f
		}
		container.Refresh()
	}}
	an.Start()
	return an
}
