package main

import (
	"fmt"
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func TestPoints(t *testing.T) {
	ts := "Points: (x:5.000 y:5.000) (x:15.000 y:5.000) (x:15.000 y:15.000) (x:5.000 y:15.000)"
	var sb strings.Builder
	c := NewMoverCircle(color.White, color.Black, 10, 10, 10, 10)
	sc := c.GetPoints().String()
	sb.WriteString(sc + "\n")
	if sc != ts {
		t.Errorf("Failed")
	}
	i := NewMoverImage(10, 10, 10, 10, canvas.NewImageFromResource(Lander_Png))
	si := i.GetPoints().String()
	sb.WriteString(si + "\n")
	if si != ts {
		t.Errorf("Failed")
	}
	r := NewMoverRect(color.White, 10, 10, 10, 10, 0)
	sr := r.GetPoints().String()
	sb.WriteString(sr + "\n")
	if sr != ts {
		t.Errorf("Failed")
	}
	app.New()

	txt := NewMoverText("HI", 100, 100, 20, fyne.TextAlignCenter)
	tr := txt.GetPoints().String()
	sb.WriteString(tr + "\n")
	if tr != "Points: (x:86.953 y:87.383) (x:113.047 y:87.383) (x:113.047 y:112.617) (x:86.953 y:112.617)" {
		t.Errorf("Failed")
	}
	fmt.Println(sb.String())
}

func TestSizeCenter(t *testing.T) {
	var sb strings.Builder
	b := NewBounds(-100, -100, 100, 100)
	sb.WriteString(testSize(t, 1, b, 200, 200))
	sb.WriteString(testCenter(t, 2, b, 0, 0))
	fmt.Println(sb.String())
}

func TestInside(t *testing.T) {
	var sb strings.Builder
	b := NewBounds(-1, -1, 1, 1)
	sb.WriteString(testInside(t, 3, b, true, -1, -1, 0, 1))
	sb.WriteString(testInside(t, 3, b, true, 0, -1, 0, 1))
	sb.WriteString(testInside(t, 3, b, true, 1, -1, 0, 1))

	sb.WriteString(testInside(t, 3, b, false, -2, -2, -1, 0, 1, 2))
	sb.WriteString(testInside(t, 3, b, false, -1, -2, 2))
	sb.WriteString(testInside(t, 3, b, false, 0, -2, 2))
	sb.WriteString(testInside(t, 3, b, false, 1, -2, 2))
	sb.WriteString(testInside(t, 3, b, false, 2, -2, -1, 0, 1, 2))

	fmt.Println(sb.String())
}

func testInside(t *testing.T, id int, b *Bounds, exp bool, p ...float64) string {
	var sb strings.Builder
	for i := 1; i < len(p); i = i + 2 {
		a := b.Inside(p[i], p[0])
		if a != exp {
			s := fmt.Sprintf("\n%d) *** Error *** INSIDE (%f, %f) %s - expected %t actual %t ", id, p[i], p[i+1], b.String(), exp, a)
			sb.WriteString(s)
			t.Error(sb.String())
		}
	}
	if len(sb.String()) > 0 {
		return sb.String()
	}
	return fmt.Sprintf("\n%d) OK: Inside: %t ", id, exp)
}

func testSize(t *testing.T, id int, b *Bounds, expw, exph float32) string {
	if b.Size().Width != expw {
		s := fmt.Sprintf("\n%d) *** Error *** Width - expected %f actual %f ", id, expw, b.Size().Width)
		t.Errorf(s)
		return s
	}
	if b.Size().Height != exph {
		s := fmt.Sprintf("\n%d) *** Error *** Height - expected %f actual %f ", id, exph, b.Size().Height)
		t.Errorf(s)
		return s
	}
	return fmt.Sprintf("\n%d) OK: Width: %f, Height: %f", id, b.Size().Width, b.Size().Height)
}
func testCenter(t *testing.T, id int, b *Bounds, expx, expy float32) string {
	if b.Center().X != expx {
		s := fmt.Sprintf("\n%d) *** Error *** Center X - expected %f actual %f ", id, expx, b.Center().X)
		t.Errorf(s)
		return s
	}
	if b.Center().Y != expy {
		s := fmt.Sprintf("\n%d) *** Error *** Center Y - expected %f actual %f ", id, expy, b.Center().Y)
		t.Errorf(s)
		return s
	}
	return fmt.Sprintf("\n%d) OK: Center X: %f, Center Y: %f", id, b.Center().X, b.Center().Y)
}
