package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type FileBrowserLineType int

const (
	FB_PARENT FileBrowserLineType = iota
	FB_DIR
	FB_FILE
	FB_ERR
	FB_UNSET
)

var (
	FB_bgColour       = color.RGBA{0, 45, 0, 0}
	FB_selectBgColour = color.RGBA{45, 0, 0, 0}
	FB_textColour     = color.White
	FB_borderColour   = color.White
)

type FileBrowserWidgetLine struct {
	visible      bool
	lineNo       int
	selectLineNo int
	yOffset      float32
	filePath     string
	lineType     FileBrowserLineType
	position     fyne.Position
	size         fyne.Size
	cText        *canvas.Text
	rect         *canvas.Rectangle
}

var _ fyne.Widget = (*FileBrowserWidgetLine)(nil)
var _ fyne.CanvasObject = (*FileBrowserWidgetLine)(nil)
var _ fyne.WidgetRenderer = (*fileBrowserWidgetLineRenderer)(nil)

func NewFileBrowserWidgetLine(text string, filePath string, lineType FileBrowserLineType, textStyle fyne.TextStyle, textSize float32, lineNo int, lineScale, width float32) *FileBrowserWidgetLine {
	t := &canvas.Text{Text: text, TextSize: textSize, TextStyle: textStyle}
	r := &canvas.Rectangle{StrokeColor: FB_borderColour, FillColor: FB_bgColour, StrokeWidth: 1}
	me := fyne.MeasureText(t.Text, t.TextSize, t.TextStyle)
	si := fyne.Size{Width: width, Height: me.Height * lineScale}
	po := fyne.Position{X: 0, Y: si.Height * float32(lineNo)}
	r.Resize(si)
	return &FileBrowserWidgetLine{cText: t, filePath: filePath, rect: r, lineNo: lineNo, selectLineNo: -1, lineType: lineType, yOffset: (si.Height - me.Height) / 2, size: si, position: po}
}

func (be *FileBrowserWidgetLine) isInside(x, y float32) bool {
	p := be.position
	s := be.size
	if x < p.X {
		return false
	}
	if y < p.Y {
		return false
	}
	if x > (p.X + s.Width) {
		return false
	}
	if y > (p.Y + s.Height) {
		return false
	}
	return true
}

// Move moves this object to the given position relative to its parent.
// This should only be called if your object is not in a container with a layout manager.
func (be *FileBrowserWidgetLine) Move(p fyne.Position) {
	be.position = p
}

// Resize resizes this object to the given size.
// This should only be called if your object is not in a container with a layout manager.
func (be *FileBrowserWidgetLine) Resize(s fyne.Size) {
	be.size = s
}

// Position returns the current position of the object relative to its parent.
func (be *FileBrowserWidgetLine) Position() fyne.Position {
	return be.position
}

// Refresh must be called if this object should be redrawn because its inner state changed.
func (be *FileBrowserWidgetLine) Refresh() {
	if be.selectLineNo != -1 {
		be.rect.FillColor = FB_selectBgColour
	} else {
		be.rect.FillColor = FB_bgColour
	}
	be.rect.Refresh()
	be.cText.Refresh()
}

// Size returns the current size of this object.
func (be *FileBrowserWidgetLine) Size() fyne.Size {
	return be.size
}

// MinSize returns the minimum size this object needs to be drawn.
func (be *FileBrowserWidgetLine) MinSize() fyne.Size {
	return fyne.Size{Width: 10, Height: 10}
}

// visibility

// Hide hides this object.
func (be *FileBrowserWidgetLine) Hide() {
	be.visible = false
}

// Visible returns whether this object is visible or not.
func (be *FileBrowserWidgetLine) Visible() bool {
	return !be.visible
}

// Show shows this object.
func (be *FileBrowserWidgetLine) Show() {
	be.visible = true
}

// Create the renderer. This is called by the fyne application
func (w *FileBrowserWidgetLine) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newFileBrowserWidgetLineRenderer(w)
}

// --------------------------------------------------------------------------------- LINE RENDERER
type fileBrowserWidgetLineRenderer struct {
	lineWidget *FileBrowserWidgetLine // Reference to the widget holding the current state
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newFileBrowserWidgetLineRenderer(w *FileBrowserWidgetLine) *fileBrowserWidgetLineRenderer {
	return &fileBrowserWidgetLineRenderer{
		lineWidget: w,
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
// Dont call r.widget.Refresh() it causes a stack overflow
func (r *fileBrowserWidgetLineRenderer) Refresh() {
	fmt.Printf("fileBrowserWidgetLineRenderer.Refresh")
}

// Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *fileBrowserWidgetLineRenderer) Layout(s fyne.Size) {
	r.lineWidget.cText.Move(fyne.Position{X: 10, Y: r.lineWidget.yOffset})
	r.lineWidget.rect.Resize(fyne.Size{Width: s.Width, Height: r.lineWidget.size.Height})
}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *fileBrowserWidgetLineRenderer) MinSize() fyne.Size {
	return r.lineWidget.MinSize()
}

// Return a list of each canvas object.
func (r *fileBrowserWidgetLineRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.lineWidget.cText, r.lineWidget.rect}
}

// Cleanup if resources have been allocated
func (r *fileBrowserWidgetLineRenderer) Destroy() {}
