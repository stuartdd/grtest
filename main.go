package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.Resize(fyne.Size{Width: 300, Height: 200})

	hello := widget.NewLabel("Hello Fyne!")
	mainWindow.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	mainWindow.ShowAndRun()
}
