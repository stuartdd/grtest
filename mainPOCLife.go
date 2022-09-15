package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

/*
-------------------------------------------------------------------- main
*/
func mainPOCLife(mainWindow fyne.Window, controller *MoverController) *fyne.Container {
	cw := controller.width
	ch := controller.height

	// gen := NewLifeGen()
	// gen.AddCell()

	cont := container.New(NewStaticLayout(cw, ch))
	return cont
}
