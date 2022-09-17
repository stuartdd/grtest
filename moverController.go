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

type AnimationController struct {
	fastAnimation *fyne.Animation
	ft            float32
	delay         int64
	tick          func(*MoverController, *AnimationController, float32)
	running       bool
}

type MoverController struct {
	size      fyne.Size
	width     float64
	height    float64
	movers    []Movable
	update    []func(float64) bool
	keyPress  func(*fyne.KeyEvent)
	animation *AnimationController
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
	c := &MoverController{size: fyne.Size{Width: float32(width), Height: float32(height)}, width: width, height: height, movers: make([]Movable, 0), update: make([]func(float64) bool, 0)}
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

func (cc *MoverController) AddOnUpdate(update func(float64) bool) {
	cc.update = append(cc.update, update)
}

func (cc *MoverController) Update(time float64) {
	if len(cc.update) > 0 {
		q := true
		for _, f := range cc.update {
			if !f(time) {
				q = false
			}
		}
		if !q {
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
	return cc.animation.running
}

func (cc *MoverController) StopAnimation() {
	if cc.animation != nil {
		cc.animation.Stop()

	}
}
func (cc *MoverController) StartAnimation() {
	if cc.animation != nil {
		cc.animation.Start()
	}
}

func controllerDefaultTick(cc *MoverController, ac *AnimationController, f float32) {
	cc.Update(float64(f - ac.ft))
	if f == 1.0 {
		ac.ft = 0
	} else {
		ac.ft = f
	}
}

func (cc *MoverController) InitAnimationController(delay int64, tick func(*MoverController, *AnimationController, float32)) {
	ac := &AnimationController{delay: delay, ft: 0, running: false}
	if tick == nil {
		tick = controllerDefaultTick
	}
	ac.tick = tick
	if ac.delay == 0 {
		aa := &fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
			ac.tick(cc, ac, f)
		}}
		ac.fastAnimation = aa
		ac.running = true
		aa.Start()
	} else {
		ac.running = false
		go func() {
			for {
				time.Sleep(time.Millisecond * time.Duration(ac.delay))
				if ac.running {
					ac.tick(cc, ac, float32(delay))
				}
			}
		}()
		ac.running = true

	}
	cc.animation = ac
}

func (ac *AnimationController) Start() {
	if ac.fastAnimation != nil {
		ac.fastAnimation.Start()
	}
	ac.running = true
}

func (ac *AnimationController) Stop() {
	if ac.fastAnimation != nil {
		ac.fastAnimation.Stop()
	}
	ac.running = false
}
