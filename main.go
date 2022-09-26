package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	mainContainer  *fyne.Container
	mainController *MoverController
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
	mainController = NewMoverController(1000, 1000)
	mainContainer = MainPOCLife(mainWindow, 1000, 1000, mainController)
	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == "Escape" {
			mainWindow.Close()
		}
		mainController.KeyPress(key)
	})

	mainWindow.SetContent(mainContainer)
	mainController.AddAfterUpdate(func(f float64) bool {
		mainContainer.Refresh()
		return true
	})

	go func() {
		time.Sleep(time.Millisecond * 500)
		mainController.InitAnimationController(mainController.GetAnimationDelay(), nil)
	}()

	mainWindow.ShowAndRun()
	mainController.StopAnimation()
}
