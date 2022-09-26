package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type FileBrowseMouseEventType int

const (
	FB_ME_NONE FileBrowseMouseEventType = 0b00000000
	FB_ME_DOWN FileBrowseMouseEventType = 0b00000001
	FB_ME_UP   FileBrowseMouseEventType = 0b00000010
	FB_ME_TAP  FileBrowseMouseEventType = 0b00000100
	FB_ME_DTAP FileBrowseMouseEventType = 0b00001000
	FB_ME_MIN  FileBrowseMouseEventType = 0b00010000
	FB_ME_MOUT FileBrowseMouseEventType = 0b00100000
	FB_ME_MOVE FileBrowseMouseEventType = 0b01000000
)

// Widget code starts here
//
// A text widget with theamed background and foreground
type FileBrowserWidget struct {
	widget.BaseWidget // Inherit from BaseWidget
	objects           []fyne.CanvasObject
	minSize           fyne.Size
	size              fyne.Size
	path              string
	onSizeChange      func(fyne.Size, fyne.Size)
	onMouseEvent      func(float32, float32, FileBrowseMouseEventType)
	onMouseMask       FileBrowseMouseEventType
}

var _ desktop.Mouseable = (*MoverWidget)(nil)
var _ fyne.Tappable = (*MoverWidget)(nil)
var _ fyne.DoubleTappable = (*MoverWidget)(nil)
var _ fyne.WidgetRenderer = (*fileBrowserWidgetRenderer)(nil)
var _ fyne.Widget = (*FileBrowserWidget)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewFileBrowserWidget(cx, cy float64, path string) *FileBrowserWidget {
	w := &FileBrowserWidget{ // Create this widget with an initial text value
		objects:     make([]fyne.CanvasObject, 0),
		minSize:     fyne.Size{Width: float32(cx), Height: float32(cy)},
		size:        fyne.Size{Width: float32(cx), Height: float32(cy)},
		onMouseMask: FB_ME_NONE,
	}
	w.ExtendBaseWidget(w) // Initialiase the BaseWidget
	w.SetPath(path, "*.rle")
	return w
}

func (w *FileBrowserWidget) SetPath(path, pattern string) {
	w.path = path
	if w.path == "" {
		return
	}
	vb := container.NewVBox()
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			vb.Add(widget.NewLabel(err.Error()))
			return err
		}
		n := info.Name()
		if !(strings.HasPrefix(n, ".") || strings.HasPrefix(path, ".") || strings.HasPrefix(n, "_") || strings.HasPrefix(path, "_")) {
			match, err := filepath.Match(pattern, info.Name())
			if match && err == nil {
				fmt.Printf("dir: %v: name: %s path: %s\n", info.IsDir(), info.Name(), path)
				vb.Add(widget.NewLabel(n))
			}
		}
		return nil
	})
	w.objects = make([]fyne.CanvasObject, 0)
	bg := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
	bg.Resize(w.size)
	w.Add(bg)
	w.Add(vb)
	if err != nil {
		fmt.Println(err)
	}
}

// Create the renderer. This is called by the fyne application
func (w *FileBrowserWidget) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newFileBrowserWidgetRenderer(w)
}

// Add canvas objects to the widget
func (w *FileBrowserWidget) Add(co fyne.CanvasObject) {
	if co == nil {
		return
	}
	w.objects = append(w.objects, co)
}

func (w *FileBrowserWidget) SetOnSizeChange(f func(fyne.Size, fyne.Size)) {
	w.onSizeChange = f
}

//
// FileBrowserWidget MOUSE HANDLING ----------------------------------------------------------- MOUSE HANDLING
//
func (mc *FileBrowserWidget) SetOnMouseEvent(f func(float32, float32, FileBrowseMouseEventType), mask FileBrowseMouseEventType) {
	mc.onMouseEvent = f
}

func (mc *FileBrowserWidget) SetOnMouseEventMask(mask FileBrowseMouseEventType) {
	if mask == FB_ME_NONE {
		mc.onMouseMask = FB_ME_NONE
	} else {
		mc.onMouseMask = mc.onMouseMask | mask
	}
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (mc *FileBrowserWidget) MouseIn(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_MIN) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_MIN)
	}
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (mc *FileBrowserWidget) MouseMoved(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_MOVE) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_MOVE)
	}
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (mc *FileBrowserWidget) MouseOut() {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_MOUT) != 0 {
		mc.onMouseEvent(0, 0, FB_ME_MOUT)
	}
}

func (mc *FileBrowserWidget) MouseDown(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_DOWN) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_DOWN)
	}
}

func (mc *FileBrowserWidget) MouseUp(me *desktop.MouseEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_UP) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_UP)
	}
}

func (mc *FileBrowserWidget) Tapped(me *fyne.PointEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_TAP) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_TAP)
	}
}

func (mc *FileBrowserWidget) DoubleTapped(me *fyne.PointEvent) {
	if mc.onMouseEvent != nil && (mc.onMouseMask&FB_ME_DTAP) != 0 {
		d := me.AbsolutePosition.X - me.Position.X
		mc.onMouseEvent(me.Position.X-d, me.Position.Y-d, FB_ME_DTAP)
	}
}

//
// RENDERER ----------------------------------------------------------------------------------- RENDERER
//
// Widget Renderer code starts here
//
type fileBrowserWidgetRenderer struct {
	widget *FileBrowserWidget // Reference to the widget holding the current state
}

//
// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newFileBrowserWidgetRenderer(myWidget *FileBrowserWidget) *fileBrowserWidgetRenderer {
	return &fileBrowserWidgetRenderer{
		widget: myWidget,
	}
}

//
// The Refresh() method is called if the state of the widget changes or the
// theme is changed
// Dont call r.widget.Refresh() it causes a stack overflow
//
func (r *fileBrowserWidgetRenderer) Refresh() {

}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *fileBrowserWidgetRenderer) Layout(s fyne.Size) {
	if r.widget.onSizeChange != nil && (r.widget.size.Width != s.Width || r.widget.size.Height != s.Height) {
		r.widget.onSizeChange(r.widget.size, s)
	}
	r.widget.size = s
}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *fileBrowserWidgetRenderer) MinSize() fyne.Size {
	return r.widget.minSize
}

// Return a list of each canvas object.
func (r *fileBrowserWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.widget.objects
}

// Cleanup if resources have been allocated
func (r *fileBrowserWidgetRenderer) Destroy() {}
