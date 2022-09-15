package main

import (
	"time"

	"fyne.io/fyne/v2"
)

type Movable interface {
	Update(float64)
	SetShouldMove(f func(Movable, float64, float64) bool)
	SetCenter(float64, float64)
	GetCenter() (float64, float64)
	SetSpeed(float64, float64)
	GetSpeed() (float64, float64)
	SetAngle(int)
	GetAngle() int
	SetAngleSpeed(float64)
	GetAngleSpeed() float64
	UpdateContainerWithObjects(*fyne.Container)
	GetSizeAndCenter() *SizeAndCenter
	GetBounds() *Bounds
	GetPoints() *Points
	SetSize(fyne.Size)
	SetVisible(bool)
	IsVisible() bool
	ContainsAny(*Points) bool
	String() string
}

type MoverController struct {
	size      fyne.Size
	width     float64
	height    float64
	movers    []Movable
	update    func(float64) bool
	keyPress  func(*fyne.KeyEvent)
	animation *fyne.Animation
}

var _ Movable = (*MoverLines)(nil)
var _ Movable = (*MoverImage)(nil)
var _ Movable = (*MoverCircle)(nil)
var _ Movable = (*MoverText)(nil)
var _ Movable = (*MoverGroup)(nil)

type StaticLayout struct {
	size fyne.Size
}

func NewStaticLayout(w, h float64) *StaticLayout {
	return &StaticLayout{size: fyne.Size{Width: float32(w), Height: float32(h)}}
}

func (sl *StaticLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return sl.size
}

func (sl *StaticLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
}

/*
-------------------------------------------------------------------- Controller
*/
func NewMoverController(width, height float64) *MoverController {
	c := &MoverController{size: fyne.Size{Width: float32(width), Height: float32(height)}, width: width, height: height, movers: make([]Movable, 0)}
	return c
}

func (cc *MoverController) KeyPress(key *fyne.KeyEvent) {
	if cc.keyPress != nil {
		cc.keyPress(key)
	}
}

func (cc *MoverController) SetOnKeyPress(keyPress func(*fyne.KeyEvent)) {
	cc.keyPress = keyPress
}

func (cc *MoverController) SetOnUpdate(update func(float64) bool) {
	cc.update = update
}

func (cc *MoverController) Update(time float64) {
	if cc.update != nil {
		if !cc.update(time) {
			return
		}
	}
	for _, m := range cc.movers {
		m.Update(time)
	}
}

func (cc *MoverController) AddMover(m Movable, c *fyne.Container) {
	cc.movers = append(cc.movers, m)
	if c == nil {
		return
	}
	m.UpdateContainerWithObjects(c)
}

func (cc *MoverController) IsAnimation() bool {
	return cc.animation != nil
}

func (cc *MoverController) StopAnimation() {
	if cc.animation != nil {
		cc.animation.Stop()
		cc.animation = nil
	}
}

func (cc *MoverController) StartAnimation() {
	cc.StopAnimation()

	var ft float32 = 0
	cc.animation = &fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		cc.Update(float64(f - ft))
		if f == 1.0 {
			ft = 0
		} else {
			ft = f
		}
	}}
	cc.animation.Start()
}
