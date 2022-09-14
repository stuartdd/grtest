package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type MoverImage struct {
	speedx  float64
	speedy  float64
	centerx float64
	centery float64

	imageSize  fyne.Size
	image      *canvas.Image
	shouldMove func(Movable, float64, float64) bool
}

/*
-------------------------------------------------------------------- MoverImage
*/
func NewMoverImage(x, y, w, h float64, image *canvas.Image) *MoverImage {
	image.Resize(fyne.Size{Width: float32(w), Height: float32(h)})
	image.FillMode = canvas.ImageFillOriginal
	im := &MoverImage{imageSize: fyne.Size{Width: float32(w), Height: float32(h)}, image: image, centerx: x, centery: y, speedx: 0, speedy: 0}
	im.SetCenter(x, y)
	return im
}

func (mv *MoverImage) String() string {
	return fmt.Sprintf("Circle x:%.3f y:%.3f w:%.3f h:%.3f", mv.centerx, mv.centery, mv.image.Size().Width, mv.image.MinSize().Height)
}

func (mv *MoverImage) ContainsAny(p *Points) bool {
	if mv.IsVisible() {
		return mv.GetBounds().ContainsAny(p)
	}
	return false
}

func (mv *MoverImage) SetVisible(v bool) {
	if !v {
		mv.image.Hide()
	} else {
		mv.image.Show()
	}
}

func (mv *MoverImage) IsVisible() bool {
	return mv.image.Visible()
}

func (mv *MoverImage) SetShouldMove(f func(Movable, float64, float64) bool) {
	mv.shouldMove = f
}

func (mv *MoverImage) Update(time float64) {
	dx := mv.speedx * time
	dy := mv.speedy * time
	if mv.shouldMove == nil || (mv.shouldMove != nil && mv.shouldMove(mv, dx, dy)) {
		mv.centerx = mv.centerx + dx
		mv.centery = mv.centery + dy
		mv.image.Move(fyne.Position{X: float32(mv.centerx) - (mv.imageSize.Width / 2), Y: float32(mv.centery) - mv.imageSize.Height/2})
	}
}

func (mv *MoverImage) GetCanvasObject() fyne.CanvasObject {
	return mv.image
}

func (mv *MoverImage) GetSizeAndCenter() *SizeAndCenter {
	return NewSizeAndCenter(float64(mv.imageSize.Width), float64(mv.imageSize.Height), mv.centerx, mv.centery)
}

func (mv *MoverImage) GetBounds() *Bounds {
	w2 := float64(mv.imageSize.Width / 2)
	h2 := float64(mv.imageSize.Height / 2)
	return &Bounds{x1: mv.centerx - w2, x2: mv.centerx + w2, y1: mv.centery - h2, y2: mv.centery + h2}
}

func (mv *MoverImage) GetPoints() *Points {
	if mv.image.Visible() {
		p := &Points{x: make([]float64, 4), y: make([]float64, 4)}
		cx := mv.centerx
		cy := mv.centery
		w2 := float64(mv.imageSize.Width / 2)
		h2 := float64(mv.imageSize.Height / 2)
		p.x[0] = cx - w2 // Top left
		p.y[0] = cy - h2
		p.x[1] = cx + w2 // Top Right
		p.y[1] = cy - h2
		p.x[2] = cx + w2 // Bottom Right
		p.y[2] = cy + h2
		p.x[3] = cx - w2 // Bottom Left
		p.y[3] = cy + h2
		return p
	}
	return &Points{x: make([]float64, 0), y: make([]float64, 0)}
}

func (mv *MoverImage) SetSize(size fyne.Size) {
	mv.image.Resize(mv.imageSize)
}

func (mv *MoverImage) SetCenter(x, y float64) {
	mv.centerx = x
	mv.centery = y
	mv.image.Move(fyne.Position{X: float32(mv.centerx) - (mv.imageSize.Width / 2), Y: float32(mv.centery) - mv.imageSize.Height/2})
}

func (mv *MoverImage) SetSpeed(x, y float64) {
	mv.speedx = x
	mv.speedy = y
}

func (mv *MoverImage) SetAngle(a int) {
}

func (mv *MoverImage) GetAngle() int {
	return 0
}

func (mv *MoverImage) SetAngleSpeed(as float64) {
}

func (mv *MoverImage) GetAngleSpeed() float64 {
	return 0
}

func (mv *MoverImage) GetCenter() (float64, float64) {
	return mv.centerx, mv.centery
}

func (mv *MoverImage) GetSpeed() (float64, float64) {
	return mv.speedx, mv.speedy
}
