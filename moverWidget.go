package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Widget code starts here
//
// A text widget with theamed background and foreground
type MoverWidget struct {
	widget.BaseWidget // Inherit from BaseWidget
	minSize           fyne.Size
	size              fyne.Size
	onSizeChange      func(fyne.Size, fyne.Size)
	bottomObjects     []fyne.CanvasObject
	topObjects        []fyne.CanvasObject
	fileWidget        fyne.Widget
	onMouseEvent      func(*MoverWidgetMouseEvent)
	onMouseMask       MoverMouseEventType
	mouseDown         bool
	mouseDownX        int64
	mouseDownY        int64
	mouseDragX        int64
	mouseDragY        int64
}

type MoverMouseEventType int

const (
	MM_ME_NONE MoverMouseEventType = 0b0000000000000000
	MM_ME_DOWN MoverMouseEventType = 0b0000000000000001
	MM_ME_UP   MoverMouseEventType = 0b0000000000000010
	MM_ME_TAP  MoverMouseEventType = 0b0000000000000100
	MM_ME_DTAP MoverMouseEventType = 0b0000000000001000
	MM_ME_MIN  MoverMouseEventType = 0b0000000000010000
	MM_ME_MOUT MoverMouseEventType = 0b0000000000100000
	MM_ME_MOVE MoverMouseEventType = 0b0000000001000000
	MM_ME_DRAG MoverMouseEventType = 0b0000000010000000
)

type MoverWidgetMouseEvent struct {
	X1       int64
	Y1       int64
	X2       int64
	Y2       int64
	Event    MoverMouseEventType
	Button   int
	Dragging bool
}

func NewMoverWidgetMouseEvent(me *desktop.MouseEvent, et MoverMouseEventType) *MoverWidgetMouseEvent {
	mwme := NewMoverWidgetMouseEventZero(et)
	mwme.Button = int(me.Button)
	d := me.AbsolutePosition.X - me.Position.X
	mwme.X1 = int64(me.Position.X - d)
	mwme.Y1 = int64(me.Position.Y - d)
	return mwme
}

func NewMoverWidgetPointEvent(me *fyne.PointEvent, et MoverMouseEventType) *MoverWidgetMouseEvent {
	mwme := NewMoverWidgetMouseEventZero(et)
	d := me.AbsolutePosition.X - me.Position.X
	mwme.X1 = int64(me.Position.X - d)
	mwme.Y1 = int64(me.Position.Y - d)
	return mwme
}

func NewMoverWidgetMouseEventZero(et MoverMouseEventType) *MoverWidgetMouseEvent {
	return &MoverWidgetMouseEvent{X1: 0, Y1: 0, X2: 0, Y2: 0, Button: 0, Event: et}
}

func (me *MoverWidgetMouseEvent) String() string {
	return fmt.Sprintf("Pos:%f Size:%f Dragging %t", me.Position(), me.Size(), me.Dragging)
}

func (me *MoverWidgetMouseEvent) Position() *fyne.Position {
	x := me.X1
	if me.X2 < x {
		x = me.X2
	}
	y := me.Y1
	if me.Y2 < y {
		y = me.Y2
	}
	return &fyne.Position{X: float32(x), Y: float32(y)}
}

func (me *MoverWidgetMouseEvent) Size() *fyne.Size {
	w := me.X2 - me.X1
	if w < 0 {
		w = w * -1
	}
	h := me.Y2 - me.Y1
	if h < 0 {
		h = h * -1
	}
	return &fyne.Size{Width: float32(w), Height: float32(h)}
}

var _ desktop.Mouseable = (*MoverWidget)(nil)
var _ fyne.Tappable = (*MoverWidget)(nil)
var _ fyne.DoubleTappable = (*MoverWidget)(nil)
var _ fyne.WidgetRenderer = (*moverWidgetRenderer)(nil)
var _ fyne.Widget = (*MoverWidget)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewMoverWidget(cx, cy float64) *MoverWidget {
	w := &MoverWidget{ // Create this widget with an initial text value
		bottomObjects: make([]fyne.CanvasObject, 0),
		topObjects:    make([]fyne.CanvasObject, 0),
		fileWidget:    nil,
		minSize:       fyne.Size{Width: float32(cx), Height: float32(cy)},
		size:          fyne.Size{Width: float32(cx), Height: float32(cy)},
		onMouseMask:   MM_ME_NONE,
	}
	w.ExtendBaseWidget(w) // Initialiase the BaseWidget
	return w
}

// Create the renderer. This is called by the fyne application
func (w *MoverWidget) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newMoverWidgetRenderer(w)
}

// Add canvas objects to the widget
func (w *MoverWidget) AddBottom(co fyne.CanvasObject) {
	if co == nil {
		return
	}
	w.bottomObjects = append(w.bottomObjects, co)
}

func (w *MoverWidget) SetFileBrowserWidget(widget fyne.Widget) {
	w.fileWidget = widget
}

// Add canvas objects to the widget
func (w *MoverWidget) AddMover(mover Movable) {
	if mover == nil {
		return
	}
	w.bottomObjects = append(w.bottomObjects, mover.GetCanvasObjects()...)
}

func (w *MoverWidget) AddTop(co fyne.CanvasObject) {
	if co == nil {
		return
	}
	w.topObjects = append(w.topObjects, co)
}

func (w *MoverWidget) SetOnSizeChange(f func(fyne.Size, fyne.Size)) {
	w.onSizeChange = f
}

func (mc *MoverWidget) SetOnMouseEvent(f func(me *MoverWidgetMouseEvent), mask MoverMouseEventType) {
	mc.onMouseEvent = f
	mc.onMouseMask = mask
}

func (mc *MoverWidget) SetOnMouseEventMask(mask MoverMouseEventType) {
	if mask == MM_ME_NONE {
		mc.onMouseMask = MM_ME_NONE
	} else {
		mc.onMouseMask = mc.onMouseMask | mask
	}
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (mc *MoverWidget) MouseIn(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_MIN) != 0 {
		mc.onMouseEvent(NewMoverWidgetMouseEvent(me, MM_ME_MIN))
	}
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (mc *MoverWidget) MouseMoved(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_MOVE) != 0 {
		mwme := NewMoverWidgetMouseEvent(me, MM_ME_MOVE)
		mwme.Dragging = mc.mouseDown
		if mc.mouseDown {
			mc.mouseDragX = mwme.X1
			mc.mouseDragY = mwme.Y1
			mwme.X2 = mc.mouseDragX
			mwme.Y2 = mc.mouseDragY
			mwme.X1 = mc.mouseDownX
			mwme.Y1 = mc.mouseDownY
		}
		mc.onMouseEvent(mwme)
	}
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (mc *MoverWidget) MouseOut() {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_MOUT) != 0 {
		mc.onMouseEvent(NewMoverWidgetMouseEventZero(MM_ME_MOUT))
	}
}

func (mc *MoverWidget) MouseDown(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_DOWN) != 0 {
		mwme := NewMoverWidgetMouseEvent(me, MM_ME_DOWN)
		mc.mouseDownX = mwme.X1
		mc.mouseDownY = mwme.Y1
		mc.mouseDragX = mwme.X1
		mc.mouseDragY = mwme.Y1
		mc.mouseDown = true
		mc.onMouseEvent(mwme)
	}
}

func (mc *MoverWidget) MouseUp(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_UP) != 0 {
		mwme := NewMoverWidgetMouseEvent(me, MM_ME_UP)
		if mc.mouseDown {
			if (mc.mouseDownX != mc.mouseDragX) || (mc.mouseDownY != mc.mouseDragY) {
				mwme.X1 = mc.mouseDownX
				mwme.Y1 = mc.mouseDownY
				mwme.X2 = mc.mouseDragX
				mwme.Y2 = mc.mouseDragY
				mwme.Event = MM_ME_DRAG
			} else {
				mwme.Event = MM_ME_TAP
			}
		}
		mc.mouseDown = false
		mc.onMouseEvent(mwme)
	}
}

// Tapped is dealt with at mouse up. This prevents a tap when ending a mouse drag
func (mc *MoverWidget) Tapped(me *fyne.PointEvent) {
}

func (mc *MoverWidget) DoubleTapped(me *fyne.PointEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&MM_ME_DTAP) != 0 {
		mc.onMouseEvent(NewMoverWidgetPointEvent(me, MM_ME_DTAP))
	}
}

// Widget Renderer code starts here
type moverWidgetRenderer struct {
	moverWidget *MoverWidget // Reference to the widget holding the current state
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newMoverWidgetRenderer(myWidget *MoverWidget) *moverWidgetRenderer {
	return &moverWidgetRenderer{
		moverWidget: myWidget,
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
// Dont call r.widget.Refresh() it causes a stack overflow
func (r *moverWidgetRenderer) Refresh() {
	if r.moverWidget.fileWidget != nil {
		r.moverWidget.fileWidget.Refresh()
	}
}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *moverWidgetRenderer) Layout(s fyne.Size) {
	if r.moverWidget.onSizeChange != nil && (r.moverWidget.size.Width != s.Width || r.moverWidget.size.Height != s.Height) {
		r.moverWidget.onSizeChange(r.moverWidget.size, s)
	}
	r.moverWidget.size = s
	if r.moverWidget.fileWidget != nil {
		r.moverWidget.fileWidget.Resize(s)
	}
}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *moverWidgetRenderer) MinSize() fyne.Size {
	return r.moverWidget.minSize
}

// Return a list of each canvas object.
func (r *moverWidgetRenderer) Objects() []fyne.CanvasObject {
	o := make([]fyne.CanvasObject, 0)
	o = append(o, moverWidget.bottomObjects...)
	o = append(o, moverWidget.topObjects...)
	if r.moverWidget.fileWidget != nil {
		o = append(o, r.moverWidget.fileWidget)
	}
	return o
}

// Cleanup if resources have been allocated
func (r *moverWidgetRenderer) Destroy() {}
