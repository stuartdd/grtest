package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

/*
-------------------------------------------------------------------- main
*/
func main() {
	a := app.New()
	mainWindow := a.NewWindow("Hello")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Close()
	})
	mainWindow.SetMaster()
	mainWindow.SetIcon(GoLogo_Png)
	controller := NewControllerContainer(500, 500)
	container := container.New(controller)
	controller.SetContainer(container)

	mainPOC(mainWindow, controller)

	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		controller.KeyPress(key)
	})
	mainWindow.SetContent(container)
	an := startAnimation(controller)
	mainWindow.ShowAndRun()
	an.Stop()
}

func startAnimation(controller *ControllerLayout) *fyne.Animation {
	var ft float32 = 0
	an := &fyne.Animation{Duration: time.Duration(time.Second), RepeatCount: 1000000, Curve: fyne.AnimationLinear, Tick: func(f float32) {
		controller.Update(float64(f - ft))
		if f == 1.0 {
			ft = 0
		} else {
			ft = f
		}
		controller.Refresh()
	}}
	an.Start()
	return an
}
