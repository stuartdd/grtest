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
	SetCenter(float64, float64)
	GetCenter() (float64, float64)
	SetSpeed(float64, float64)
	GetSpeed() (float64, float64)
	GetCanvasObject() fyne.CanvasObject
	GetSizeAndCenter() *SizeAndCenter
	GetBounds() *Bounds
	SetSize(fyne.Size)
}

type SizeAndCenter struct {
	Width   float64
	Height  float64
	CenterX float64
	CenterY float64
}

type Bounds struct {
	x1 float64
	y1 float64
	x2 float64
	y2 float64
}

type MoverCircle struct {
	speedx  float64
	speedy  float64
	centerx float64
	centery float64
	width   float64
	height  float64
	circle  *canvas.Circle
}

type MoverLines struct {
	speedx     float64
	accx       float64
	speedy     float64
	accy       float64
	currentAng int
	accAng     float64
	speedAng   float64
	centerx    float64
	centery    float64
	lines      []*canvas.Line
}

type MoverImage struct {
	speedx  float64
	speedy  float64
	centerx float64
	centery float64

	imageSize fyne.Size
	image     *canvas.Image
}

type ControllerLayout struct {
	size   fyne.Size
	movers []Movable
}

var _ Movable = (*MoverLines)(nil)
var _ Movable = (*MoverImage)(nil)
var _ Movable = (*MoverCircle)(nil)

/*
-------------------------------------------------------------------- SizeAndCenter
*/
func NewSizeAndCenter(w, h, x, y float64) *SizeAndCenter {
	return &SizeAndCenter{Width: w, Height: h, CenterX: x, CenterY: y}
}

func (sc *SizeAndCenter) Size() fyne.Size {
	return fyne.Size{Height: float32(sc.Height), Width: float32(sc.Width)}
}

func (sc *SizeAndCenter) Center() fyne.Position {
	return fyne.Position{X: float32(sc.CenterX), Y: float32(sc.CenterY)}
}

/*
-------------------------------------------------------------------- MoverCircle
*/
func NewMoverCircle(strokeColour, fillColour color.Color, centerx, centery, width, height float64) *MoverCircle {
	circle := &canvas.Circle{
		Position1:   fyne.Position{X: float32(centerx - width/2), Y: float32(centery - height/2)},
		Position2:   fyne.Position{X: float32(centerx + width/2), Y: float32(centery + height/2)},
		StrokeColor: strokeColour, FillColor: fillColour, StrokeWidth: 1,
	}
	return &MoverCircle{speedx: 0, speedy: 0, centerx: centerx, centery: centery, width: width, height: height, circle: circle}
}

func (mv *MoverCircle) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if (dx != 0) || (dy != 0) {
		mv.centerx = mv.centerx + dx
		mv.centery = mv.centery + dy
		mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
		mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
		mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
		mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
	}
}

func (mv *MoverCircle) GetCanvasObject() fyne.CanvasObject {
	return mv.circle
}

func (mv *MoverCircle) GetSizeAndCenter() *SizeAndCenter {
	return NewSizeAndCenter(mv.width, mv.height, mv.centerx, mv.centery)
}

func (mv *MoverCircle) GetBounds() *Bounds {
	w2 := float64(mv.width / 2)
	h2 := float64(mv.height / 2)
	return &Bounds{x1: mv.centerx - w2, x2: mv.centerx + w2, y1: mv.centery - h2, y2: mv.centery + h2}
}

func (mv *MoverCircle) SetSize(size fyne.Size) {
	mv.width = float64(size.Width)
	mv.height = float64(size.Height)
	mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
	mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
	mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
	mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
}

func (mv *MoverCircle) Init() {
}

func (mv *MoverCircle) SetCenter(x, y float64) {
	mv.centerx = x
	mv.centery = y
	mv.circle.Position1.X = float32(mv.centerx - (mv.width / 2))
	mv.circle.Position1.Y = float32(mv.centery - (mv.height / 2))
	mv.circle.Position2.X = float32(mv.centerx + (mv.width / 2))
	mv.circle.Position2.Y = float32(mv.centery + (mv.height / 2))
}

func (mv *MoverCircle) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverCircle) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
}

func (mv *MoverCircle) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

/*
-------------------------------------------------------------------- MoverLines
*/

func NewMoverLines(centerx, centery, speedAng float64) *MoverLines {
	return &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerx: centerx, centery: centery, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
}

func NewMoverRect(colour color.Color, centerx, centery, w, h, speedAng float64) *MoverLines {
	ml := &MoverLines{speedx: 0, speedy: 0, speedAng: speedAng, centerx: centerx, centery: centery, currentAng: 0, accAng: 0, lines: make([]*canvas.Line, 0)}
	ml.AddLine(colour, float32(centerx-w/2), float32(centery-h/2), float32(centerx+w/2), float32(centery-h/2))
	ml.AddLineToo(colour, float32(centerx+w/2), float32(centery+h/2))
	ml.AddLineToo(colour, float32(centerx-w/2), float32(centery+h/2))
	ml.AddLineToo(colour, float32(centerx-w/2), float32(centery-h/2))
	return ml
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
	for mv.accAng > 360.0 {
		mv.accAng = mv.accAng - 360.0
	}
	for mv.accAng < 0 {
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
			rotatePosition(mv.centerx, mv.centery, &l.Position1, ra)
			rotatePosition(mv.centerx, mv.centery, &l.Position2, ra)
		}
	}
	mv.centerx = mv.centerx + float64(dx)
	mv.centery = mv.centery + float64(dy)
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
	currentSize := mv.GetSizeAndCenter()
	scaleX := float64(size.Width) / currentSize.Width
	scaleY := float64(size.Height) / currentSize.Height
	for _, l := range mv.lines {
		scalePoint(mv.centerx, mv.centery, &l.Position1, scaleX, scaleY)
		scalePoint(mv.centerx, mv.centery, &l.Position2, scaleX, scaleY)
	}
}

func (mv *MoverLines) GetSizeAndCenter() *SizeAndCenter {
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
	w := float64(maxx - minx)
	h := float64(maxy - miny)
	return NewSizeAndCenter(w, h, (float64(minx) + w/2), (float64(miny) + h/2))
}

func (mv *MoverLines) GetBounds() *Bounds {
	sac := mv.GetSizeAndCenter()
	w2 := sac.Width / 2
	h2 := sac.Height / 2
	return &Bounds{x1: sac.CenterX - w2, x2: sac.CenterX + w2, y1: sac.CenterY - h2, y2: sac.CenterY + h2}
}

func (mv *MoverLines) Init() {
}

func (mv *MoverLines) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverLines) SetCenter(x, y float64) {
	dx := x - mv.centerx
	dy := y - mv.centery
	for _, l := range mv.lines {
		l.Position1.X = l.Position1.X + float32(dx)
		l.Position2.X = l.Position2.X + float32(dx)
		l.Position1.Y = l.Position1.Y + float32(dy)
		l.Position2.Y = l.Position2.Y + float32(dy)
	}
	mv.centerx = x
	mv.centery = y
}

func (mv *MoverLines) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
}

func (mv *MoverLines) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}

func (mv *MoverLines) AddLine(colour color.Color, x1, y1, x2, y2 float32) {
	mv.lines = append(mv.lines, &canvas.Line{StrokeColor: colour, StrokeWidth: 1, Position1: fyne.Position{X: x1, Y: y1}, Position2: fyne.Position{X: x2, Y: y2}})
}

func (mv *MoverLines) AddLines(colour color.Color, xy ...float32) {
	for i := 0; i < len(xy); i = i + 2 {
		mv.AddLineToo(colour, xy[i], xy[i+1])
	}
}

func (mv *MoverLines) AddLineToo(colour color.Color, x2, y2 float32) {
	var x1 float32 = 0.0
	var y1 float32 = 0.0
	le := len(mv.lines)
	if le > 0 {
		lf := mv.lines[le-1]
		x1 = lf.Position2.X
		y1 = lf.Position2.Y
	}
	mv.AddLine(colour, x1, y1, x2, y2)
}

/*
-------------------------------------------------------------------- MoverImage
*/
func NewMoverImage(x, y, w, h float64, image *canvas.Image) *MoverImage {
	return &MoverImage{imageSize: fyne.Size{Width: float32(w), Height: float32(h)}, image: image, centerx: x, centery: y, speedx: 0, speedy: 0}
}

func (mv *MoverImage) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if (dx != 0) || (dy != 0) {
		mv.centerx = mv.centerx + dx
		mv.centery = mv.centery + dy
		mv.image.Move(fyne.Position{X: float32(mv.centerx) - (mv.imageSize.Width / 2), Y: float32(mv.centery) - mv.imageSize.Height/2})
	}
}

func (mv *MoverImage) GetCanvasObject() fyne.CanvasObject {
	return mv.image
}

func (mv *MoverImage) GetSizeAndCenter() *SizeAndCenter {
	return NewSizeAndCenter(float64(mv.imageSize.Width), float64(mv.imageSize.Height), mv.centerx, mv.centery)
}

func (mv *MoverImage) GetBounds() *Bounds {
	w2 := float64(mv.imageSize.Width / 2)
	h2 := float64(mv.imageSize.Height / 2)
	return &Bounds{x1: mv.centerx - w2, x2: mv.centerx + w2, y1: mv.centery - h2, y2: mv.centery + h2}
}

func (mv *MoverImage) SetSize(size fyne.Size) {
	mv.image.Resize(mv.imageSize)
}

func (mv *MoverImage) Init() {
	mv.image.Resize(mv.imageSize)
	mv.image.FillMode = canvas.ImageFillOriginal
}

func (mv *MoverImage) SetCenter(x, y float64) {
	mv.centerx = x
	mv.centery = y
}

func (mv *MoverImage) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverImage) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
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
	moverImage2.SetSpeed(10, 10)

	lines1 := NewMoverLines(200, 200, 10)
	lines1.AddLine(color.White, 200, 200, 200, 300)
	lines1.AddLineToo(color.White, 100, 300)
	lines1.AddLineToo(color.White, 200, 200)
	lines1.SetSpeed(2, 2)

	bBox1 := NewMoverRect(color.RGBA{250, 0, 0, 255}, 200, 200, 100, 100, 0)
	bBox2 := NewMoverRect(color.RGBA{0, 0, 255, 255}, 200, 200, 100, 100, 0)

	lines3 := NewMoverLines(0, 0, 0)
	lines3.AddLineToo(color.White, 1000, 1000)

	circ1 := NewMoverCircle(color.RGBA{255, 255, 0, 255}, color.RGBA{0, 255, 255, 255}, 400, 100, 100, 100)
	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Close()
	})
	mainWindow.SetMaster()
	mainWindow.SetIcon(GoLogo_Png)

	controller := NewControllerContainer(500, 500)
	container := container.New(controller)

	controller.Add(circ1)
	container.Add(circ1.GetCanvasObject())
	controller.Add(moverImage2)
	container.Add(moverImage2.GetCanvasObject())

	controller.Add(lines1)
	container.Add(lines1.GetCanvasObject())
	controller.Add(lines3)
	container.Add(lines3.GetCanvasObject())

	container.Add(bBox1.GetCanvasObject())
	container.Add(bBox2.GetCanvasObject())

	mainWindow.SetContent(container)
	controller.Init()

	an := startAnimation(controller, container)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			SetSpeedAndTarget(circ1, moverImage2, 12)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			s := lines1.GetSizeAndCenter()
			bBox1.SetSize(s.Size())
			bBox1.SetCenter(float64(s.CenterX), float64(s.CenterY))
			s = moverImage2.GetSizeAndCenter()
			bBox2.SetSize(s.Size())
			bBox2.SetCenter(float64(s.CenterX), float64(s.CenterY))
			scaleMovable(circ1, 0.998)
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
