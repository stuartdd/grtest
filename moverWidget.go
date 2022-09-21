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
type MyWidget struct {
	widget.BaseWidget // Inherit from BaseWidget
	size              fyne.Size
	objects           []fyne.CanvasObject
}

func (mc *MyWidget) MouseDown(*desktop.MouseEvent) {
	fmt.Println("MouseDown")
}

func (mc *MyWidget) MouseUp(*desktop.MouseEvent) {
	fmt.Println("MouseUp")
}

func (mc *MyWidget) Tapped(*fyne.PointEvent) {
	fmt.Println("Tapped")
}

var _ desktop.Mouseable = (*MyWidget)(nil)
var _ fyne.Tappable = (*MyWidget)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewMyWidget(cx, cy float64) *MyWidget {
	w := &MyWidget{ // Create this widget with an initial text value
		objects: make([]fyne.CanvasObject, 0),
		size:    fyne.Size{Width: float32(cx), Height: float32(cy)},
	}
	w.ExtendBaseWidget(w) // Initialiase the BaseWidget
	return w
}

// Create the renderer. This is called by the fyne application
func (w *MyWidget) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newMyWidgetRenderer(w)
}

// Create the renderer. This is called by the fyne application
func (w *MyWidget) Add(co fyne.CanvasObject) {
	// Pass this widget to the renderer so it can access the text field
	w.objects = append(w.objects, co)
}

// Widget Renderer code starts here
type myWidgetRenderer struct {
	widget *MyWidget // Reference to the widget holding the current state
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newMyWidgetRenderer(myWidget *MyWidget) *myWidgetRenderer {
	return &myWidgetRenderer{
		widget: myWidget,
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
//
// Note: The background and foreground colours are set from the current theme
func (r *myWidgetRenderer) Refresh() {
}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *myWidgetRenderer) Layout(s fyne.Size) {
	r.widget.size = s
}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *myWidgetRenderer) MinSize() fyne.Size {
	return r.widget.size
}

// Return a list of each canvas object.
func (r *myWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.widget.objects
}

// Cleanup if resources have been allocated
func (r *myWidgetRenderer) Destroy() {}
