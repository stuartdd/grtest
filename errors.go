package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var ERR_TEXT_STYLE = fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2}

type FixedErrorLayout struct {
	w float32
	h float32
}

type ErrorContainer struct {
	callback     func()
	prefix       string
	errorText    *canvas.Text
	actionButton *widget.Button
	background   *canvas.Rectangle
	container    *fyne.Container
}

func NewErrorContainer(prefix, errorText, action string, callback func()) *ErrorContainer {
	ec := &ErrorContainer{callback: callback, prefix: prefix}
	ec.actionButton = widget.NewButton(action, func() {
		if ec.callback == nil {
			ec.Close()
		} else {
			ec.callback()
		}
	})
	ec.background = &canvas.Rectangle{FillColor: color.RGBA{200, 0, 0, 0}}
	ec.errorText = &canvas.Text{TextStyle: ERR_TEXT_STYLE, TextSize: 19, Color: color.Black}
	ec.container = container.New(newErrorLayout(100, 30), ec.background, ec.errorText, ec.actionButton)
	ec.SetErrorString(errorText)
	return ec
}

func (ec *ErrorContainer) ShowError(err error) {
}

func (ec *ErrorContainer) Container() *fyne.Container {
	return ec.container
}

func (ec *ErrorContainer) Close() {
	ec.errorText.Text = ""
	ec.container.Hide()
}

func (ec *ErrorContainer) SetErrorString(err string) {
	if err != "" {
		ec.errorText.Text = fmt.Sprintf(" %s : %s", ec.prefix, err)
		ec.container.Show()
	} else {
		ec.errorText.Text = ""
		ec.container.Hide()
	}
}

func newErrorLayout(w float32, h float32) *FixedErrorLayout {
	return &FixedErrorLayout{w: w, h: h}
}

func (d *FixedErrorLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(d.w, d.h)
}

func (d *FixedErrorLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var bw float32 = 0
	for _, o := range objects {
		_, ok := o.(*widget.Button)
		if ok {
			bw = o.MinSize().Width
		}
	}

	for _, o := range objects {
		_, ok := o.(*canvas.Rectangle)
		if ok {
			o.Resize(fyne.Size{Width: containerSize.Width, Height: d.h})
			o.Move(fyne.Position{X: 0, Y: 0})
		}
		_, ok = o.(*canvas.Text)
		if ok {
			o.Resize(fyne.Size{Width: containerSize.Width - bw, Height: d.h})
			o.Move(fyne.Position{X: 0, Y: 0})
		}
		_, ok = o.(*widget.Button)
		if ok {
			o.Move(fyne.Position{X: containerSize.Width - bw, Y: 0})
			o.Resize(fyne.Size{Width: bw, Height: d.h})

		}
	}
}

var _ fyne.Layout = (*FixedErrorLayout)(nil)
