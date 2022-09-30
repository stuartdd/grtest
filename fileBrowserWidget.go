package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
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
	textStyle         *fyne.TextStyle
	textSize          float32
	currentPath       string
	pattern           string
	onMouseEvent      func(float32, float32, FileBrowseMouseEventType)
	onMouseMask       FileBrowseMouseEventType
	err               error
}

var _ desktop.Mouseable = (*MoverWidget)(nil)
var _ fyne.Tappable = (*MoverWidget)(nil)
var _ fyne.DoubleTappable = (*MoverWidget)(nil)
var _ fyne.WidgetRenderer = (*fileBrowserWidgetRenderer)(nil)
var _ fyne.Widget = (*FileBrowserWidget)(nil)
var _ fyne.CanvasObject = (*FileBrowserWidget)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewFileBrowserWidget(cx, cy float64, pattern string) *FileBrowserWidget {
	w := &FileBrowserWidget{ // Create this widget with an initial text value
		objects:     make([]fyne.CanvasObject, 0),
		minSize:     fyne.Size{Width: float32(cx), Height: float32(cy)},
		size:        fyne.Size{Width: float32(cx), Height: float32(cy)},
		textStyle:   &fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2},
		textSize:    fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText),
		onMouseMask: FB_ME_NONE,
		currentPath: "",
		pattern:     "*",
		err:         nil,
	}
	w.ExtendBaseWidget(w) // Initialiase the BaseWidget
	w.BaseWidget.Resize(w.size)
	return w
}

func (w *FileBrowserWidget) GetSelected() (string, FileBrowserLineType) {
	for _, o := range w.objects {
		ow, ok := o.(*FileBrowserWidgetLine)
		if ok {
			if ow.selectLineNo >= 0 {
				return path.Join(w.currentPath, ow.cText.Text), ow.lineType
			}
		}
	}
	return "", FB_ERR
}

func (w *FileBrowserWidget) SetParentPath() {
	pp, err := PathToParentPath(w.currentPath)
	if err == nil {
		w.SetPath(pp, w.pattern)
	}
}

func (w *FileBrowserWidget) SetPath(newPath, pattern string) {
	if newPath == "" {
		return
	}
	line := 0
	coFile := make([]fyne.CanvasObject, 0)
	coDir := make([]fyne.CanvasObject, 0)
	_, err := PathToParentPath(newPath)
	if err == nil {
		fbe := NewFileBrowserWidgetLine(".. (up to parent directory)", FB_PARENT, *w.textStyle, w.textSize, line, 2, w.size.Width)
		coDir = append(coDir, fbe)
		line++
	}
	files, err := os.ReadDir(newPath)
	if err != nil {
		fbe := NewFileBrowserWidgetLine(err.Error(), FB_ERR, *w.textStyle, w.textSize, line, 2, w.size.Width)
		coDir = append(coDir, fbe)
		w.objects = coDir
		return
	}
	if len(files) == 0 {
		fbe := NewFileBrowserWidgetLine("No Files Found", FB_ERR, *w.textStyle, w.textSize, line, 2, w.size.Width)
		coDir = append(coDir, fbe)
		w.objects = coDir
		return
	}
	for _, info := range files {
		n := info.Name()
		if !(strings.HasPrefix(n, ".") || strings.HasPrefix(n, "_")) {
			if info.IsDir() {
				fbe := NewFileBrowserWidgetLine(n, FB_DIR, *w.textStyle, w.textSize, line, 2, w.size.Width)
				coDir = append(coDir, fbe)
				line++
			} else {
				match, err := filepath.Match(pattern, info.Name())
				if match && err == nil {
					fbe := NewFileBrowserWidgetLine(n, FB_FILE, *w.textStyle, w.textSize, line, 2, w.size.Width)
					coFile = append(coFile, fbe)
					line++
				}
			}
		}
	}
	if len(files) == 0 {
		fbe := NewFileBrowserWidgetLine("No matching files returned", FB_ERR, *w.textStyle, w.textSize, line, 2, w.size.Width)
		coDir = append(coDir, fbe)
		w.objects = coDir
		return
	}
	w.pattern = pattern
	w.currentPath = newPath
	w.objects = append(coDir, coFile...)
}

func (mc *FileBrowserWidget) SelectByMouse(x, y float32) int {
	lin := -1
	for _, o := range mc.objects {
		fbwl, ok := o.(*FileBrowserWidgetLine)
		if ok {
			if fbwl.isInside(x, y) {
				lin = fbwl.lineNo
				fbwl.selectLineNo = lin
			} else {
				fbwl.selectLineNo = -1
			}
		}
	}
	return lin
}

// Create the renderer. This is called by the fyne application
func (w *FileBrowserWidget) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newFileBrowserWidgetRenderer(w)
}

// FileBrowserWidget MOUSE HANDLING ----------------------------------------------------------- MOUSE HANDLING
func (mc *FileBrowserWidget) SetOnMouseEvent(f func(float32, float32, FileBrowseMouseEventType), mask FileBrowseMouseEventType) {
	mc.onMouseEvent = f
	mc.onMouseMask = mask
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

// RENDERER ----------------------------------------------------------------------------------- RENDERER
//
// Widget Renderer code starts here
type fileBrowserWidgetRenderer struct {
	widget *FileBrowserWidget // Reference to the widget holding the current state
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newFileBrowserWidgetRenderer(myWidget *FileBrowserWidget) *fileBrowserWidgetRenderer {
	return &fileBrowserWidgetRenderer{
		widget: myWidget,
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
// Dont call r.widget.Refresh() it causes a stack overflow
func (r *fileBrowserWidgetRenderer) Refresh() {
	objects := r.widget.objects
	for _, o := range objects {
		o.Refresh()
	}
}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *fileBrowserWidgetRenderer) Layout(s fyne.Size) {
	objects := r.widget.objects
	for i, o := range objects {
		o.Resize(fyne.Size{Width: s.Width, Height: o.Size().Height})
		o.Move(fyne.Position{X: 0, Y: float32(i) * o.Size().Height})
	}
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
