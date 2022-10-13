package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	err               error
	saveLabel         *widget.Label
	saveEntry         *widget.Entry
	saveCancel        *widget.Button
	saveCommit        *widget.Button
	saveForm          *fyne.Container
	saveError         *widget.Label

	onSelectedEvent  func(string, string) error
	onFileFoundEvent func(fs.DirEntry, string, FileBrowserLineType) string
	onSaveEvent      func(string, bool, error) error
}

var _ desktop.Mouseable = (*MoverWidget)(nil)
var _ fyne.Tappable = (*MoverWidget)(nil)
var _ fyne.DoubleTappable = (*MoverWidget)(nil)
var _ fyne.WidgetRenderer = (*fileBrowserWidgetRenderer)(nil)
var _ fyne.Widget = (*FileBrowserWidget)(nil)
var _ fyne.CanvasObject = (*FileBrowserWidget)(nil)
var FBW_TEXT_STYLE = fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2}

type FixedWHLayout struct {
	w float32
	h float32
}

func NewFixedWHLayout(w float32, h float32) *FixedWHLayout {
	return &FixedWHLayout{w: w, h: h}
}

func (d *FixedWHLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(d.w, d.h)
}

func (d *FixedWHLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var wid float32 = 0
	for _, o := range objects {
		_, ok := o.(*widget.Entry)
		if !ok {
			wid = wid + o.MinSize().Width
		}
	}
	var pos float32 = 0
	var min float32 = 0
	for _, o := range objects {
		_, ok := o.(*widget.Entry)
		if ok {
			min = containerSize.Width - wid
		} else {
			min = o.MinSize().Width
		}
		o.Resize(fyne.NewSize(min, d.h))
		o.Move(fyne.Position{X: pos, Y: 0})
		pos = pos + min
	}
}

var _ fyne.Layout = (*FixedWHLayout)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewFileBrowserWidget(cx, cy float64) *FileBrowserWidget {
	w := &FileBrowserWidget{ // Create this widget with an initial text value
		objects:     make([]fyne.CanvasObject, 0),
		minSize:     fyne.Size{Width: float32(cx), Height: float32(cy)},
		size:        fyne.Size{Width: float32(cx), Height: float32(cy)},
		textStyle:   &fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2},
		textSize:    fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText),
		currentPath: "",
		err:         nil,
	}
	w.ExtendBaseWidget(w) // Initialiase the BaseWidget
	w.BaseWidget.Resize(w.size)
	return w
}

func (w *FileBrowserWidget) saveActionNotify(save bool, err error) {
	if w.onSaveEvent != nil {
		ent := strings.TrimSpace(w.saveEntry.Text)
		if ent == "" && save {
			return
		}
		err := w.onSaveEvent(path.Join(w.currentPath, ent), save, err)
		if err != nil {
			w.saveEntry.Hide()
			w.saveError.Text = fmt.Sprintf(" ERROR:%s", err.Error())
			w.saveError.Show()
		} else {
			w.saveEntry.Show()
			w.saveError.Text = ""
			w.saveError.Hide()
		}
	}
}
func (w *FileBrowserWidget) saveDataChanged(s string) {
	if w.onSaveEvent != nil {
		ent := strings.TrimSpace(w.saveEntry.Text)
		if ent == "" {
			return
		}
		full := path.Join(w.currentPath, ent)
		_, errStat := os.Stat(full)
		if os.IsNotExist(errStat) {
			w.saveCommit.Text = "Save"
		} else {
			w.saveCommit.Text = "* Overwrite *"
		}
		w.saveForm.Refresh()
	}
}

func (mc *FileBrowserWidget) updateFileEntry(s string) {
	x := mc.saveEntry.OnChanged
	mc.saveEntry.OnChanged = nil
	mc.saveEntry.Text = s
	mc.saveEntry.OnChanged = x
	mc.saveDataChanged(s)
}

func (w *FileBrowserWidget) InputSaveForm(prompt string) *fyne.Container {
	w.saveLabel = widget.NewLabelWithStyle(prompt, fyne.TextAlignCenter, FBW_TEXT_STYLE)
	w.saveEntry = widget.NewEntry()
	w.saveEntry.TextStyle = FBW_TEXT_STYLE
	w.saveEntry.PlaceHolder = "Enter a file name"
	w.saveCancel = widget.NewButton("Cancel", func() {
		w.saveActionNotify(false, nil)
	})
	w.saveCommit = widget.NewButton("Save", func() {
		w.saveActionNotify(true, nil)
	})
	w.saveCommit.Importance = widget.HighImportance
	w.saveEntry = widget.NewEntry()
	w.saveError = widget.NewLabelWithStyle("", fyne.TextAlignCenter, FBW_TEXT_STYLE)
	w.saveEntry.OnSubmitted = func(s string) {
		w.saveActionNotify(true, nil)
	}
	w.saveEntry.OnChanged = w.saveDataChanged
	w.saveError.Hide()
	w.saveForm = container.New(NewFixedWHLayout(100, 40), w.saveLabel, w.saveError, w.saveEntry, w.saveCancel, w.saveCommit)
	return w.saveForm
}

func (w *FileBrowserWidget) GetSelected() (string, FileBrowserLineType) {
	for _, o := range w.objects {
		ow, ok := o.(*FileBrowserWidgetLine)
		if ok {
			if ow.selected {
				return ow.filePath, ow.lineType
			}
		}
	}
	return "", FB_ERR
}

func (w *FileBrowserWidget) SetParentPath() {
	pp, err := PathToParentPath(w.currentPath)
	if err == nil {
		w.SetPath(pp)
	}
}

func (mc *FileBrowserWidget) SetPath(newPath string) {
	if newPath == "" {
		return
	}
	line := 0
	coFile := make([]fyne.CanvasObject, 0)
	coDir := make([]fyne.CanvasObject, 0)
	_, err := PathToParentPath(newPath)
	if err == nil && newPath != "/" {
		fbe := NewFileBrowserWidgetLine(".. (up to parent directory)", "..", FB_PARENT, *mc.textStyle, mc.textSize, line, 2, mc.size.Width)
		coDir = append(coDir, fbe)
		line++
	}
	files, err := os.ReadDir(newPath)
	if err != nil {
		fbe := NewFileBrowserWidgetLine(err.Error(), "", FB_ERR, *mc.textStyle, mc.textSize, line, 2, mc.size.Width)
		coDir = append(coDir, fbe)
		mc.objects = coDir
		return
	}
	if len(files) == 0 {
		fbe := NewFileBrowserWidgetLine("No Files Found", "", FB_ERR, *mc.textStyle, mc.textSize, line, 2, mc.size.Width)
		coDir = append(coDir, fbe)
		mc.objects = coDir
		return
	}
	for _, info := range files {
		typ := FB_FILE
		if info.IsDir() {
			typ = FB_DIR
		}
		n := info.Name()
		if mc.onFileFoundEvent != nil {
			n = mc.onFileFoundEvent(info, newPath, typ)
		}
		if n != "" {
			fbe := NewFileBrowserWidgetLine(n, path.Join(newPath, info.Name()), typ, *mc.textStyle, mc.textSize, line, 2, mc.size.Width)
			if typ == FB_FILE {
				coFile = append(coFile, fbe)
			} else {
				coDir = append(coDir, fbe)
			}
			line++
		}
	}
	mc.currentPath = newPath
	if mc.saveForm != nil {
		mc.saveDataChanged(mc.saveEntry.Text)
	}
	mc.objects = append(coDir, coFile...)
}

func (mc *FileBrowserWidget) selectByMouse(x, y float32) *FileBrowserWidgetLine {
	var fbwlSel *FileBrowserWidgetLine
	for _, o := range mc.objects {
		fbwl, ok := o.(*FileBrowserWidgetLine)
		if ok {
			if fbwl.isInside(x, y) {
				fbwl.selected = true
				fbwlSel = fbwl
			} else {
				fbwl.selected = false
			}
		}
	}
	return fbwlSel
}

// Create the renderer. This is called by the fyne application
func (mc *FileBrowserWidget) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newFileBrowserWidgetRenderer(mc)
}

// FileBrowserWidget MOUSE HANDLING ----------------------------------------------------------- MOUSE HANDLING
func (mc *FileBrowserWidget) SetOnFileFoundEvent(f func(fs.DirEntry, string, FileBrowserLineType) string) {
	mc.onFileFoundEvent = f
}

func (mc *FileBrowserWidget) SetOnSelectedEvent(f func(string, string) error) {
	mc.onSelectedEvent = f
	mc.onSaveEvent = nil
}

func (mc *FileBrowserWidget) SetOnSaveEvent(f func(string, bool, error) error) {
	mc.onSaveEvent = f
	mc.onSelectedEvent = nil
	mc.saveEntry.OnSubmitted = func(s string) {
		mc.saveActionNotify(true, nil)
	}
}

func (mc *FileBrowserWidget) Tapped(me *fyne.PointEvent) {
	if mc.onSelectedEvent != nil {
		d := me.AbsolutePosition.X - me.Position.X
		fbwl := mc.selectByMouse(me.Position.X-d, me.Position.Y-d)
		if fbwl != nil {
			if mc.saveForm != nil && fbwl.lineType == FB_FILE {
				_, f := path.Split(fbwl.filePath)
				mc.updateFileEntry(f)
			}
			go mc.Refresh()
		}
	}
}

func (mc *FileBrowserWidget) DoubleTapped(me *fyne.PointEvent) {
	d := me.AbsolutePosition.X - me.Position.X
	fbwl := mc.selectByMouse(me.Position.X-d, me.Position.Y-d)
	if fbwl != nil {
		switch fbwl.lineType {
		case FB_PARENT:
			go mc.SetParentPath()
		case FB_DIR:
			go fbWidget.SetPath(fbwl.filePath)
		case FB_FILE:
			if mc.onSelectedEvent == nil {
				_, f := path.Split(fbwl.filePath)
				mc.updateFileEntry(f)
			} else {
				err := mc.onSelectedEvent(fbwl.filePath, mc.currentPath)
				if err == nil {
					mc.Hide()
				}
			}
		}
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
