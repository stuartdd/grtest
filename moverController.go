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
	GetCanvasObjects() []fyne.CanvasObject
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
	movers       []Movable
	beforeUpdate []func(float64) bool
	afterUpdate  []func(float64) bool
	keyPress     func(*fyne.KeyEvent)
	animation    *AnimationController
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
	c := &MoverController{movers: make([]Movable, 0), beforeUpdate: make([]func(float64) bool, 0)}
	return c
}

func (mc *MoverController) KeyPress(key *fyne.KeyEvent) {
	if mc.keyPress != nil {
		mc.keyPress(key)
	}
}

func (mc *MoverController) SetOnKeyPress(keyPress func(*fyne.KeyEvent)) {
	mc.keyPress = keyPress
}

func (mc *MoverController) AddBeforeUpdate(update func(float64) bool) {
	mc.beforeUpdate = append(mc.beforeUpdate, update)
}

func (mc *MoverController) AddAfterUpdate(update func(float64) bool) {
	mc.afterUpdate = append(mc.afterUpdate, update)
}

// Main update fllo for animation of background loop
//
// If any updateBefore returns false then no movers are updated and no updateAfters are called
// Each movers update is called
// If an updateAfter return false then the following updateAfters and not called.
func (mc *MoverController) Update(time float64) {
	if len(mc.beforeUpdate) > 0 {
		q := true
		for _, f := range mc.beforeUpdate {
			if !f(time) {
				q = false
			}
		}
		if !q {
			return
		}
	}

	for _, m := range mc.movers {
		m.Update(time)
	}

	if len(mc.afterUpdate) > 0 {
		for _, f := range mc.afterUpdate {
			if !f(time) {
				return
			}
		}
	}
}

func (mc *MoverController) AddMover(m Movable) {
	mc.movers = append(mc.movers, m)
}

func (mc *MoverController) IsAnimation() bool {
	return mc.animation.running
}

func (mc *MoverController) StopAnimation() {
	if mc.animation != nil {
		mc.animation.Stop()

	}
}
func (mc *MoverController) StartAnimation() {
	if mc.animation != nil {
		mc.animation.Start()
	}
}

func controllerDefaultTick(mc *MoverController, ac *AnimationController, f float32) {
	mc.Update(float64(f - ac.ft))
	if f == 1.0 {
		ac.ft = 0
	} else {
		ac.ft = f
	}
}

func (mc *MoverController) InitAnimationController(delay int64, tick func(*MoverController, *AnimationController, float32)) {
	ac := &AnimationController{delay: delay, ft: 0, running: false}
	if tick == nil {
		tick = controllerDefaultTick
	}
	ac.tick = tick
	if ac.delay == 0 {
		aa := &fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
			ac.tick(mc, ac, f)
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
					ac.tick(mc, ac, float32(ac.delay))
				}
			}
		}()
		ac.running = true

	}
	mc.animation = ac
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
