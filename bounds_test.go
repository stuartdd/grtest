package main

import (
	"fmt"
	"strings"
	"testing"
)

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
