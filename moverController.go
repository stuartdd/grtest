package main

import (
	"sync"

	"fyne.io/fyne/v2"
)

type Movable interface {
	Update(float64)
	SetShouldMove(f func(float64, float64, float64, float64) bool)
	SetCenter(float64, float64)
	GetCenter() (float64, float64)
	SetSpeed(float64, float64)
	GetSpeed() (float64, float64)
	SetAngle(int)
	GetAngle() int
	SetAngleSpeed(float64)
	GetAngleSpeed() float64
	GetCanvasObject() fyne.CanvasObject
	GetSizeAndCenter() *SizeAndCenter
	GetBounds() *Bounds
	GetPoints() *Points
	SetSize(fyne.Size)
	SetVisible(bool)
	IsVisible() bool
	ContainsAny(*Points) bool
	String() string
}

type ControllerLayout struct {
	size      fyne.Size
	movers    []Movable
	container *fyne.Container
	update    func(float64) bool
	keyPress  func(*fyne.KeyEvent)
	mu        sync.Mutex
}

var _ Movable = (*MoverLines)(nil)
var _ Movable = (*MoverImage)(nil)
var _ Movable = (*MoverCircle)(nil)
var _ Movable = (*MoverText)(nil)
var _ Movable = (*MoverGroup)(nil)

/*
-------------------------------------------------------------------- ControllerLayout
*/
func NewControllerContainer(width, height float32) *ControllerLayout {
	c := &ControllerLayout{size: fyne.Size{Width: width, Height: height}, movers: make([]Movable, 0)}
	return c
}

func (cc *ControllerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	cc.size = size
}

func (cc *ControllerLayout) KeyPress(key *fyne.KeyEvent) {
	if cc.keyPress != nil {
		cc.keyPress(key)
	}
}

func (cc *ControllerLayout) SetOnKeyPress(keyPress func(*fyne.KeyEvent)) {
	cc.keyPress = keyPress
}

func (cc *ControllerLayout) SetOnUpdate(update func(float64) bool) {
	cc.update = update
}

func (cc *ControllerLayout) Update(time float64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cc.update != nil {
		if !cc.update(time) {
			return
		}
	}
	for _, m := range cc.movers {
		m.Update(time)
	}
}

func (cc *ControllerLayout) AddToContainer(co fyne.CanvasObject) {
	if cc.container == nil {
		panic("ControllerLayout requires a fyne.Container. Please use SetContainer(c)")
	}
	cc.container.Add(co)
}

func (cc *ControllerLayout) AddMover(m Movable) {
	if cc.container == nil {
		panic("ControllerLayout requires a fyne.Container. Please use SetContainer(c)")
	}
	cc.movers = append(cc.movers, m)
	cc.container.Add(m.GetCanvasObject())
}

func (cc *ControllerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return cc.size
}

func (cc *ControllerLayout) GetContainer() *fyne.Container {
	return cc.container
}

func (cc *ControllerLayout) SetContainer(container *fyne.Container) {
	cc.container = container
}

func (cc *ControllerLayout) Refresh() {
	cc.container.Refresh()
}
