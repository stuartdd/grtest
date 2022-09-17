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
	mainContainer = mainPOCLife(mainWindow, mainController)
	desktop
	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		// fmt.Println(key.Name)
		if key.Name == "Escape" {
			mainWindow.Close()
		}
		if key.Name == "F1" {
			if mainController.IsAnimation() {
				mainController.StopAnimation()
			} else {
				mainController.StartAnimation()
			}

		}
		mainController.KeyPress(key)
	})
	mainWindow.SetContent(mainContainer)
	mainController.AddOnUpdate(func(f float64) bool {
		mainContainer.Refresh()
		return false
	})

	go func() {
		time.Sleep(time.Millisecond * 500)
		mainController.InitAnimationController(50, nil)
		// for {
		// 	time.Sleep(time.Millisecond * 200)
		// 	mainContainer.Refresh()
		// }
	}()

	mainWindow.ShowAndRun()
	mainController.StopAnimation()
}
