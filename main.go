package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"
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
	GetCanvasObject() fyne.CanvasObject
}

type MoverLines struct {
	speedx   float32
	speedy   float32
	speedRot float32
	centerX  float64
	centerY  float64
	lines    []*canvas.Line
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

var _ Movable = (*MoverLines)(nil)
var _ Movable = (*MoverImage)(nil)

/*
-------------------------------------------------------------------- MoverLines
*/

func NewMoverLines(centerX, centerY float64, speedRot float32) *MoverLines {
	return &MoverLines{speedx: 0, speedy: 0, speedRot: speedRot, centerX: centerX, centerY: centerY, lines: make([]*canvas.Line, 0)}
}

func (mv *MoverLines) Update(time float32) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	ra := mv.speedRot * (time * 10)
	for _, l := range mv.lines {
		if dx != 0 {
			l.Position1.X = l.Position1.X + dx
			l.Position2.X = l.Position2.X + dx
		}
		if dy != 0 {
			l.Position1.Y = l.Position1.Y + dy
			l.Position2.Y = l.Position2.Y + dy
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

func (mv *MoverLines) AdjustSpeed(x, y float32) {
	mv.speedx = mv.speedx + x
	mv.speedy = mv.speedy + y
}

func (mv *MoverLines) GetSpeed() (float32, float32) {
	return mv.speedx, mv.speedy
}

func (mv *MoverLines) AddLine(x1, y1, x2, y2 float32, colour color.Color) {
	line := canvas.NewLine(colour)
	line.Position1.X = x1
	line.Position1.Y = y1
	line.Position2.X = x2
	line.Position2.Y = y2
	mv.lines = append(mv.lines, line)
}

func (mv *MoverLines) AddLineToo(x2, y2 float32, colour color.Color) {
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

func rotatePoint(centerX, centerY float64, point *fyne.Position, radians float32) {
	angle := int(radians * (180 / math.Pi))
	px := float64(point.X) - centerX
	py := float64(point.Y) - centerY
	point.X = float32(cos(angle)*px - sin(angle)*py + centerX)
	point.Y = float32(sin(angle)*px + cos(angle)*py + centerY)
}

/*
-------------------------------------------------------------------- MoverImage
*/
func NewMoverImage(x, y, w, h float32, image *canvas.Image) *MoverImage {
	return &MoverImage{imageSize: fyne.Size{Width: w, Height: h}, image: image, posx: x, posy: y, speedx: 0, speedy: 0}
}

func (mv *MoverImage) Update(time float32) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if (dx != 0) || (dy != 0) {
		mv.posx = mv.posx + dx
		mv.posy = mv.posy + dy
		mv.image.Move(fyne.Position{X: mv.posx, Y: mv.posy})
	}
}

func (mv *MoverImage) GetCanvasObject() fyne.CanvasObject {
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
	// testCos()
	// testSine()
	// os.Exit(0)
	moverImage1 := NewMoverImage(0, 0, 40, 40, canvas.NewImageFromResource(Lander_Png))
	moverImage1.AdjustSpeed(100, 100)
	moverImage2 := NewMoverImage(0, 0, 20, 20, canvas.NewImageFromResource(Lander_Png))
	moverImage2.AdjustSpeed(5, 5)

	lines1 := NewMoverLines(150, 150, 0.01)
	lines1.AddLine(100, 100, 150, 150, color.White)
	lines1.AddLineToo(200, 100, color.White)
	lines1.AdjustSpeed(5, 5)

	lines2 := NewMoverLines(0, 0, 0)
	lines2.AddLineToo(400, 400, color.White)

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

	mainWindow.SetContent(container)
	controller.Init()

	var ft float32 = 0
	an := fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		controller.Update(f - ft)
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

var SIN_COS_TABLE = []float64{0.000000000000, 0.017452406437, 0.034899496703, 0.052335956243, 0.069756473744, 0.087155742748, 0.104528463268, 0.121869343405, 0.139173100960, 0.156434465040, 0.173648177667, 0.190808995377, 0.207911690818, 0.224951054344, 0.241921895600, 0.258819045103, 0.275637355817, 0.292371704723, 0.309016994375, 0.325568154457, 0.342020143326, 0.358367949545, 0.374606593416, 0.390731128489, 0.406736643076, 0.422618261741, 0.438371146789, 0.453990499740, 0.469471562786, 0.484809620246, 0.500000000000, 0.515038074910, 0.529919264233, 0.544639035015, 0.559192903471, 0.573576436351, 0.587785252292, 0.601815023152, 0.615661475326, 0.629320391050, 0.642787609687, 0.656059028991, 0.669130606359, 0.681998360062, 0.694658370459, 0.707106781187, 0.719339800339, 0.731353701619, 0.743144825477, 0.754709580223, 0.766044443119, 0.777145961457, 0.788010753607, 0.798635510047, 0.809016994375, 0.819152044289, 0.829037572555, 0.838670567945, 0.848048096156, 0.857167300702, 0.866025403784, 0.874619707139, 0.882947592859, 0.891006524188, 0.898794046299, 0.906307787037, 0.913545457643, 0.920504853452, 0.927183854567, 0.933580426497, 0.939692620786, 0.945518575599, 0.951056516295, 0.956304755963, 0.961261695938, 0.965925826289, 0.970295726276, 0.974370064785, 0.978147600734, 0.981627183448, 0.984807753012, 0.987688340595, 0.990268068742, 0.992546151641, 0.994521895368, 0.996194698092, 0.997564050260, 0.998629534755, 0.999390827019, 0.999847695156, 1.000000000000, 0.999847695156, 0.999390827019, 0.998629534755, 0.997564050260, 0.996194698092, 0.994521895368, 0.992546151641, 0.990268068742, 0.987688340595, 0.984807753012, 0.981627183448, 0.978147600734, 0.974370064785, 0.970295726276, 0.965925826289, 0.961261695938, 0.956304755963, 0.951056516295, 0.945518575599, 0.939692620786, 0.933580426497, 0.927183854567, 0.920504853452, 0.913545457643, 0.906307787037, 0.898794046299, 0.891006524188, 0.882947592859, 0.874619707139, 0.866025403784, 0.857167300702, 0.848048096156, 0.838670567945, 0.829037572555, 0.819152044289, 0.809016994375, 0.798635510047, 0.788010753607, 0.777145961457, 0.766044443119, 0.754709580223, 0.743144825477, 0.731353701619, 0.719339800339, 0.707106781187, 0.694658370459, 0.681998360062, 0.669130606359, 0.656059028991, 0.642787609687, 0.629320391050, 0.615661475326, 0.601815023152, 0.587785252292, 0.573576436351, 0.559192903471, 0.544639035015, 0.529919264233, 0.515038074910, 0.500000000000, 0.484809620246, 0.469471562786, 0.453990499740, 0.438371146789, 0.422618261741, 0.406736643076, 0.390731128489, 0.374606593416, 0.358367949545, 0.342020143326, 0.325568154457, 0.309016994375, 0.292371704723, 0.275637355817, 0.258819045103, 0.241921895600, 0.224951054344, 0.207911690818, 0.190808995377, 0.173648177667, 0.156434465040, 0.139173100960, 0.121869343405, 0.104528463268, 0.087155742748, 0.069756473744, 0.052335956243, 0.034899496703, 0.017452406437}

func testSine() {

	fmt.Println(len(SIN_COS_TABLE))

	var sb1 strings.Builder
	var sb2 strings.Builder

	for i := 0; i < 360; i++ {
		ra := float64(i) * (math.Pi / 180)
		co := math.Sin(ra)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") {
			so = so[1:]
		}
		sb1.WriteString(so)
	}
	fmt.Println(sb1.String())
	fmt.Println("------------------------")
	for i := 0; i < 360; i++ {
		co := sin(i)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") {
			so = so[1:]
		}
		sb2.WriteString(so)
	}
	fmt.Println(sb2.String())
	fmt.Println("------------------------")

	b1 := []byte(sb1.String())
	b2 := []byte(sb2.String())
	if len(b1) != len(b2) {
		fmt.Printf("Len 1 %d. Len 2 %d", len(b1), len(b2))
	}
	var sb3 strings.Builder
	for i := 0; i < len(b1); i++ {
		sb3.WriteByte(b1[i])
		if b1[i] != b2[i] {
			fmt.Printf("Char at b1[%d]=%c b2[%d]=%c\n%s", i, b1[i], i, b2[i], sb3.String())
			break
		}
	}
}

func testCos() {

	fmt.Println(len(SIN_COS_TABLE))

	var sb1 strings.Builder
	var sb2 strings.Builder

	for i := 0; i < 360; i++ {
		ra := float64(i) * (math.Pi / 180)
		co := math.Cos(ra)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") {
			so = so[1:]
		}
		sb1.WriteString(so)
	}
	fmt.Println(sb1.String())
	fmt.Println("------------------------")
	for i := 0; i < 360; i++ {
		co := cos(i)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") {
			so = so[1:]
		}
		sb2.WriteString(so)
	}
	fmt.Println(sb2.String())
	fmt.Println("------------------------")

	b1 := []byte(sb1.String())
	b2 := []byte(sb2.String())
	if len(b1) != len(b2) {
		fmt.Printf("Len 1 %d. Len 2 %d", len(b1), len(b2))
	}
	var sb3 strings.Builder
	for i := 0; i < len(b1); i++ {
		sb3.WriteByte(b1[i])
		if b1[i] != b2[i] {
			fmt.Printf("Char at b1[%d]=%c b2[%d]=%c\n%s", i, b1[i], i, b2[i], sb3.String())
			break
		}
	}
}

func sin(deg int) float64 {
	deg = deg % 360
	if deg >= 180 {
		return SIN_COS_TABLE[deg%180] * -1
	}
	return SIN_COS_TABLE[deg]
}

func cos(deg int) float64 {
	deg = (deg + 90) % 360
	if deg >= 180 {
		return SIN_COS_TABLE[deg%180] * -1
	}
	return SIN_COS_TABLE[deg]
}
