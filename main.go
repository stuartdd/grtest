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
	mainContainer = mainPOClots(mainWindow, mainController)

	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		// fmt.Println(key.Name)
		if key.Name == "Escape" {
			mainWindow.Close()
		}
		if key.Name == "F1" {
			if mainController.GetAnimation() {
				mainController.StopAnimation()
			} else {
				mainController.StartAnimation()
			}

		}
		mainController.KeyPress(key)
	})

	mainWindow.SetContent(mainContainer)

	go func() {
		time.Sleep(time.Millisecond * 500)
		mainController.StartAnimation()
		for {
			time.Sleep(time.Millisecond * 50)
			mainContainer.Refresh()
		}
	}()

	mainWindow.ShowAndRun()
	mainController.StopAnimation()
}
