package main

import (
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
	objects           []fyne.CanvasObject
	onMouseEvent      func(float32, float32, MoverMouseEventType)
	onMouseMask       MoverMouseEventType
}

type MoverMouseEventType int

const (
	ME_NONE MoverMouseEventType = 0b00000000
	ME_DOWN MoverMouseEventType = 0b00000001
	ME_UP   MoverMouseEventType = 0b00000010
	ME_TAP  MoverMouseEventType = 0b00000100
	ME_DTAP MoverMouseEventType = 0b00001000
	ME_MIN  MoverMouseEventType = 0b00010000
	ME_MOUT MoverMouseEventType = 0b00100000
	ME_MOVE MoverMouseEventType = 0b01000000
)

func (mc *MoverWidget) SetOnMouseEvent(f func(float32, float32, MoverMouseEventType), mask MoverMouseEventType) {
	mc.onMouseEvent = f
	mc.SetOnMouseEventMask(mask)
}

func (mc *MoverWidget) SetOnMouseEventMask(mask MoverMouseEventType) {
	if mask == ME_NONE {
		mc.onMouseMask = ME_NONE
	} else {
		mc.onMouseMask = mc.onMouseMask | mask
	}
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (mc *MoverWidget) MouseIn(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_MIN) != 0 {
		mc.onMouseEvent(me.Position.X, me.Position.Y, ME_MIN)
	}
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (mc *MoverWidget) MouseMoved(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_MOVE) != 0 {
		mc.onMouseEvent(me.Position.X, me.Position.Y, ME_MOVE)
	}
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (mc *MoverWidget) MouseOut() {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_MOUT) != 0 {
		mc.onMouseEvent(0, 0, ME_MOUT)
	}
}

func (mc *MoverWidget) MouseDown(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_DOWN) != 0 {
		mc.onMouseEvent(me.Position.X, me.Position.Y, ME_DOWN)
	}
}

func (mc *MoverWidget) MouseUp(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_UP) != 0 {
		mc.onMouseEvent(me.Position.X, me.Position.Y, ME_UP)
	}
}

func (mc *MoverWidget) Tapped(me *fyne.PointEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_TAP) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, ME_TAP)
	}
}

func (mc *MoverWidget) DoubleTapped(me *fyne.PointEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&ME_DTAP) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, ME_DTAP)
	}
}

var _ desktop.Mouseable = (*MoverWidget)(nil)
var _ fyne.Tappable = (*MoverWidget)(nil)
var _ fyne.DoubleTappable = (*MoverWidget)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewMoverWidget(cx, cy float64) *MoverWidget {
	w := &MoverWidget{ // Create this widget with an initial text value
		objects:     make([]fyne.CanvasObject, 0),
		minSize:     fyne.Size{Width: float32(cx), Height: float32(cy)},
		size:        fyne.Size{Width: float32(cx), Height: float32(cy)},
		onMouseMask: ME_NONE,
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
func (w *MoverWidget) Add(co fyne.CanvasObject) {
	if co == nil {
		return
	}
	w.objects = append(w.objects, co)
}

// Add canvas objects to the widget
func (w *MoverWidget) AddMover(mover Movable) {
	if mover == nil {
		return
	}
	w.objects = append(w.objects, mover.GetCanvasObjects()...)
}

func (w *MoverWidget) SetOnSizeChange(f func(fyne.Size, fyne.Size)) {
	w.onSizeChange = f
}

// Widget Renderer code starts here
type moverWidgetRenderer struct {
	widget *MoverWidget // Reference to the widget holding the current state
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newMoverWidgetRenderer(myWidget *MoverWidget) *moverWidgetRenderer {
	return &moverWidgetRenderer{
		widget: myWidget,
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
//
// Note: The background and foreground colours are set from the current theme
func (r *moverWidgetRenderer) Refresh() {
}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *moverWidgetRenderer) Layout(s fyne.Size) {
	if r.widget.onSizeChange != nil && (r.widget.size.Width != s.Width || r.widget.size.Height != s.Height) {
		r.widget.onSizeChange(r.widget.size, s)
	}
	r.widget.size = s
}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *moverWidgetRenderer) MinSize() fyne.Size {
	return r.widget.minSize
}

// Return a list of each canvas object.
func (r *moverWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.widget.objects
}

// Cleanup if resources have been allocated
func (r *moverWidgetRenderer) Destroy() {}
