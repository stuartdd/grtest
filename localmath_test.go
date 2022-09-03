package main

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
)

func TestScale(t *testing.T) {
	var sb strings.Builder
	pos := &fyne.Position{X: 100, Y: 0}
	scalePoint(0, 0, pos, 0.5, 0.5)
	sb.WriteString(testPos(t, 1, 50.0, 0, pos))
	scalePoint(0, 0, pos, 100/50, 0)
	sb.WriteString(testPos(t, 2, 100.0, 0, pos))

	pos = &fyne.Position{X: 50, Y: 50}
	scalePoint(100, 100, pos, 0.5, 0.5)
	sb.WriteString(testPos(t, 3, 75.0, 75.0, pos))
	scalePoint(100, 100, pos, 50/25, 50/25)
	sb.WriteString(testPos(t, 4, 50.0, 50.0, pos))
	t.Error(sb.String())
}

func TestRotate(t *testing.T) {
	var sb strings.Builder
	pos := &fyne.Position{X: 100, Y: 0}
	sb.WriteString(testPos(t, 1, 100.0, 0, pos))
	rotatePoint(0, 0, pos, 90)
	sb.WriteString(testPos(t, 2, 0.0, 100, pos))
	rotatePoint(0, 0, pos, -90)
	sb.WriteString(testPos(t, 3, 100.0, 0, pos))
	rotatePoint(0, 0, pos, -90)
	sb.WriteString(testPos(t, 4, 0.0, -100, pos))
	rotatePoint(0, 0, pos, 45)
	rotatePoint(0, 0, pos, 45)
	sb.WriteString(testPos(t, 5, 100.0, 0, pos))
	for i := 0; i < 90; i++ {
		rotatePoint(0, 0, pos, 1)
	}
	sb.WriteString(testPos(t, 6, 0.0, 100, pos))
	for i := 0; i < 90; i++ {
		rotatePoint(0, 0, pos, 1)
	}
	sb.WriteString(testPos(t, 7, -100.0, 0, pos))
	for i := 0; i < 90; i++ {
		rotatePoint(0, 0, pos, 1)
	}
	sb.WriteString(testPos(t, 8, 0, -100, pos))
	for i := 0; i < 90; i++ {
		rotatePoint(0, 0, pos, 1)
	}
	sb.WriteString(testPos(t, 9, 100, 0, pos))
	pos = &fyne.Position{X: 100, Y: 0}
	sb.WriteString(testPos(t, 1, 100.0, 0, pos))
	for i := 0; i < 180; i++ {
		rotatePoint(0, 0, pos, -1)
	}
	sb.WriteString(testPos(t, 10, -100.0, 0, pos))
	for i := 0; i < 180; i++ {
		rotatePoint(0, 0, pos, -1)
	}
	sb.WriteString(testPos(t, 11, 100.0, 0, pos))

	fmt.Println(sb.String())
}

func testPos(t *testing.T, id int, x, y float64, p *fyne.Position) string {
	var sb strings.Builder
	spX := fmt.Sprintf("%.2f", p.X)
	if spX == "-0.00" {
		spX = spX[1:]
	}
	sX := fmt.Sprintf("%.2f", x)
	spY := fmt.Sprintf("%.2f", p.Y)
	if spY == "-0.00" {
		spY = spY[1:]
	}
	sY := fmt.Sprintf("%.2f", y)
	if spX != sX {
		s := fmt.Sprintf("\n%d) Error in X value X:%s != %s Y:%s", id, spX, sX, spY)
		sb.WriteString(s)
		t.Errorf(s)
	}
	if spY != sY {
		s := fmt.Sprintf("\n%d) Error in Y value X:%s Y:%s != %s", id, spX, spY, sY)
		sb.WriteString(s)
		t.Errorf(s)
	}
	if len(sb.String()) == 0 {
		sb.WriteString(fmt.Sprintf("\n%d) X:%s (%s) Y:%s (%s)", id, spX, sX, spY, sY))
	}
	return sb.String()
}

func TestSpeedSin(t *testing.T) {
	// Warm up the test.
	for i := 0; i < 10; i++ {
		speedTestSin()
	}
	// Run the test.
	et1, et2 := speedTestSin()
	fact := float32(et1) / float32(et2)
	// if sin faster than Sin then fact must be > 1
	if fact < 1 {
		t.Errorf("Sin() %d sin() %d factor %f should be at least 1 if sin() is faster", et1, et2, fact)
	}
	fmt.Printf("Speed factor for Sin is %f\n", fact)
}

func TestSpeedCos(t *testing.T) {
	// Warm up the test.
	for i := 0; i < 10; i++ {
		speedTestCos()
	}
	// Run the test.
	et1, et2 := speedTestCos()
	fact := float32(et1) / float32(et2)
	// if cos faster than Cos then fact must be > 1
	if fact < 1 {
		t.Errorf("Cos() %d cos() %d factor %f should be at least 1 if sin() is faster", et1, et2, fact)
	}
	fmt.Printf("Speed factor for Cos is %f\n", fact)
}

func speedTestSin() (int, int) {
	st1 := time.Now().Nanosecond()
	for i := -730; i < 730; i++ {
		ra1 := float64(i) * (math.Pi / 180)
		math.Sin(ra1)
	}
	et1 := time.Now().Nanosecond() - st1

	st2 := time.Now().Nanosecond()
	for i := -730; i < 730; i++ {
		sin(i)
	}
	et2 := time.Now().Nanosecond() - st2
	return et1, et2
}

func speedTestCos() (int, int) {
	st1 := time.Now().Nanosecond()
	for i := -730; i < 730; i++ {
		ra1 := float64(i) * (math.Pi / 180)
		math.Cos(ra1)
	}
	et1 := time.Now().Nanosecond() - st1

	st2 := time.Now().Nanosecond()
	for i := -730; i < 730; i++ {
		cos(i)
	}
	et2 := time.Now().Nanosecond() - st2
	return et1, et2
}

func TestFullRangeSin(t *testing.T) {
	for i := -730; i < 730; i++ {
		ra1 := float64(i) * (math.Pi / 180)
		co1 := math.Sin(ra1)
		so1 := fmt.Sprintf("%0.12f, ", co1)
		if strings.HasPrefix(so1, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so1 = so1[1:]
		}
		co2 := sin(i)
		so2 := fmt.Sprintf("%0.12f, ", co2)
		if strings.HasPrefix(so2, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so2 = so2[1:]
		}
		if so1 != so2 {
			t.Errorf("Sin(%d)=%s NOT sin(%d)=%s", i, so1, i, so2)
		}
	}
}

func TestFullRangeCos(t *testing.T) {
	for i := -730; i < 730; i++ {
		ra1 := float64(i) * (math.Pi / 180)
		co1 := math.Cos(ra1)
		so1 := fmt.Sprintf("%0.12f, ", co1)
		if strings.HasPrefix(so1, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so1 = so1[1:]
		}
		co2 := cos(i)
		so2 := fmt.Sprintf("%0.12f, ", co2)
		if strings.HasPrefix(so2, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so2 = so2[1:]
		}
		if so1 != so2 {
			t.Errorf("Cos(%d)=%s NOT cos(%d)=%s", i, so1, i, so2)
		}
	}
}

func TestValuesTable(t *testing.T) {
	if len(SIN_COS_TABLE) != 180 {
		t.Errorf("SIN_COS_TABLE must be 180 entries")
	}
	if SIN_COS_TABLE[0] != 0 {
		t.Errorf("SIN_COS_TABLE[0] must be 0.000000000000")
	}

	if SIN_COS_TABLE[90] != 1.000000000000 {
		t.Errorf("SIN_COS_TABLE[90] must be 1.000000000000")
	}

	// First half values values 0..90 must increase
	for i := 0; i < 91; i++ {
		if SIN_COS_TABLE[i] < 0 {
			t.Errorf("SIN_COS_TABLE[%d] must be positive", i)
		}
		if i > 0 {
			if SIN_COS_TABLE[i-1] >= SIN_COS_TABLE[i] {
				t.Errorf("SIN_COS_TABLE[%d]=%f must be less than SIN_COS_TABLE[%d]=%f", i-1, SIN_COS_TABLE[i-1], i, SIN_COS_TABLE[i])
			}
		}
	}
	// Second half values values 91..180 must decrease
	for i := 91; i < 180; i++ {
		if SIN_COS_TABLE[i] < 0 {
			t.Errorf("SIN_COS_TABLE[%d] must be positive", i)
		}
		if i > 0 {
			if SIN_COS_TABLE[i-1] <= SIN_COS_TABLE[i] {
				t.Errorf("SIN_COS_TABLE[%d]=%f must be less than SIN_COS_TABLE[%d]=%f", i-1, SIN_COS_TABLE[i-1], i, SIN_COS_TABLE[i])
			}
		}
	}
}

/*
Tests that the sin(deg) function produces the same value of the math.Sin(rad) function
*/
func TestSinFunction(t *testing.T) {
	var sb1 strings.Builder
	var sb2 strings.Builder

	for i := 0; i < 360; i++ {
		ra := float64(i) * (math.Pi / 180)
		co := math.Sin(ra)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so = so[1:]
		}
		sb1.WriteString(so)
	}
	fmt.Println(sb1.String())
	fmt.Println("------------------------")
	for i := 0; i < 360; i++ {
		co := sin(i)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so = so[1:]
		}
		sb2.WriteString(so)
	}
	fmt.Println(sb2.String())
	fmt.Println("------------------------")

	b1 := []byte(sb1.String())
	b2 := []byte(sb2.String())
	if len(b1) != len(b2) {
		t.Errorf("Len 1 %d. Len 2 %d", len(b1), len(b2))
	}
	var sb3 strings.Builder
	for i := 0; i < len(b1); i++ {
		sb3.WriteByte(b1[i])
		if b1[i] != b2[i] {
			t.Errorf("Char at b1[%d]=%c b2[%d]=%c\n%s", i, b1[i], i, b2[i], sb3.String())
			break
		}
	}
}

/*
Tests that the cos(deg) function produces the same value of the math.Cos(rad) function
*/
func TestCosFunction(t *testing.T) {
	var sb1 strings.Builder
	var sb2 strings.Builder

	for i := 0; i < 360; i++ {
		ra := float64(i) * (math.Pi / 180)
		co := math.Cos(ra)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so = so[1:]
		}
		sb1.WriteString(so)
	}
	fmt.Println(sb1.String())
	fmt.Println("------------------------")
	for i := 0; i < 360; i++ {
		co := cos(i)
		so := fmt.Sprintf("%0.12f, ", co)
		if strings.HasPrefix(so, "-0.000000000000") { // Fix the -0 returned by some values. It can be ignored in actual usage but it breaks the test.
			so = so[1:]
		}
		sb2.WriteString(so)
	}
	fmt.Println(sb2.String())
	fmt.Println("------------------------")

	b1 := []byte(sb1.String())
	b2 := []byte(sb2.String())
	if len(b1) != len(b2) {
		t.Errorf("Len 1 %d. Len 2 %d", len(b1), len(b2))
	}
	var sb3 strings.Builder
	for i := 0; i < len(b1); i++ {
		sb3.WriteByte(b1[i])
		if b1[i] != b2[i] {
			t.Errorf("Char at b1[%d]=%c b2[%d]=%c\n%s", i, b1[i], i, b2[i], sb3.String())
			break
		}
	}
}
