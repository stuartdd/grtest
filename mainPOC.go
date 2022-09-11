package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var (
	player    Movable
	textStyle = fyne.TextStyle{Bold: false, Italic: false, Monospace: true, Symbol: false, TabWidth: 2}
)

func keyPressed(key *fyne.KeyEvent) {
	sx, sy := player.GetCenter()
	switch key.Name {
	case fyne.KeyDown:
		player.SetCenter(sx, sy+1)
	case fyne.KeyUp:
		player.SetCenter(sx, sy-1)
	case fyne.KeyLeft:
		player.SetCenter(sx-1, sy)
	case fyne.KeyRight:
		player.SetCenter(sx+1, sy)
	}
}

/*
-------------------------------------------------------------------- main
*/
func mainPOC(mainWindow fyne.Window, controller *ControllerLayout) {

	player = NewMoverImage(100, 100, 40, 40, canvas.NewImageFromResource(Lander_Png))
	player.SetSpeed(15, 15)

	lines1 := NewMoverLines(200, 200, 10)
	lines1.AddLine(color.White, 200, 200, 200, 300)
	lines1.AddLineToo(color.White, 100, 300)
	lines1.AddLineToo(color.White, 200, 200)
	lines1.SetSpeed(0, 0)
	text1 := NewMoverText("Center  :", 100, 300, 10, fyne.TextAlignCenter)
	text1.SetSpeed(5, 5)
	group1 := NewMoverCroup(lines1)
	setPoints(group1, lines1)
	group1.Add(text1)
	group1.SetAngleSpeed(10)

	guideLine := NewMoverLines(0, 0, 0)
	guideLine.AddLineToo(color.White, 1000, 1000)

	text2 := NewMoverText("Trailing:", 200, 40, 20, fyne.TextAlignTrailing)
	text3 := NewMoverText("Leading :", 200, 70, 20, fyne.TextAlignLeading)
	bBox2 := NewMoverRect(color.RGBA{250, 0, 0, 255}, 200, 200, 100, 100, 0)
	bBox3 := NewMoverRect(color.RGBA{250, 0, 0, 255}, 200, 200, 100, 100, 0)
	bBox4 := NewMoverRect(color.RGBA{0, 255, 0, 255}, 200, 200, 100, 100, 0)
	bBox5 := NewMoverRect(color.RGBA{0, 255, 0, 255}, 200, 200, 100, 100, 0)

	circ1 := NewMoverCircle(color.RGBA{255, 255, 0, 255}, color.RGBA{0, 255, 255, 255}, 400, 100, 20, 20)

	/*
		Add Movers that are managed directly by the controller
	*/
	controller.AddMover(circ1)
	controller.AddMover(player)
	controller.AddMover(group1)
	controller.AddMover(text2)
	controller.AddMover(text3)
	/*
		Add Movers that are *NOT* managed directly by the controller
	*/
	controller.AddToContainer(guideLine.GetCanvasObject())
	controller.AddToContainer(bBox2.GetCanvasObject())
	controller.AddToContainer(bBox3.GetCanvasObject())
	controller.AddToContainer(bBox4.GetCanvasObject())
	controller.AddToContainer(bBox5.GetCanvasObject())

	controller.SetOnKeyPress(func(key *fyne.KeyEvent) {
		keyPressed(key)
	})

	go func() {
		for {
			time.Sleep(time.Second)
			if player.IsVisible() {
				SetSpeedAndTarget(circ1, player, 23)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 5)
			text1.SetText("Center  : MISSED")
			text2.SetText("Trailing: MISSED")
			text3.SetText("Leading : MISSED")
			text2.SetSpeed(0, 0)
			text3.SetSpeed(0, 0)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			i := player.GetBounds()
			bBox2.SetSize(i.Size())
			bBox2.SetCenter(float64(i.Center().X), float64(i.Center().Y))
			t1 := text1.GetBounds()
			bBox3.SetSize(t1.Size())
			bBox3.SetCenter(float64(t1.Center().X), float64(t1.Center().Y))
			t2 := text2.GetBounds()
			bBox4.SetSize(t2.Size())
			bBox4.SetCenter(float64(t2.Center().X), float64(t2.Center().Y))
			t3 := text3.GetBounds()
			bBox5.SetSize(t3.Size())
			bBox5.SetCenter(float64(t3.Center().X), float64(t3.Center().Y))
			if circ1.ContainsAny(player.GetPoints()) {
				text1.SetText("Center  : HIT")
				text2.SetText("Trailing: HIT")
				text3.SetText("Leading : HIT")
				text2.SetSpeed(10, 10)
				text3.SetSpeed(10, 10)
				player.SetSpeed(0, 0)
				player.SetVisible(false)
			}
		}
	}()
}

func setPoints(group *MoverGroup, mv Movable) {
	mp := mv.GetPoints()
	for i := 0; i < len(mp.x); i++ {
		co := NewMoverCircle(color.White, color.RGBA{250, 0, 0, 255}, mp.x[i], mp.y[i], 5, 5)
		group.Add(co)
	}
}