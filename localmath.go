package main

import (
	"math"

	"fyne.io/fyne/v2"
)

const (
	maxFloat32 float32 = math.MaxFloat32
	minFloat32 float32 = -math.MaxFloat32
	radToDeg   float64 = (180 / math.Pi)
)

/*
180 values generated via the test function TestSinFunction. Copy the first 180 entries (ONLY positive values) from the console output into this array
*/
var SIN_COS_TABLE = []float64{0.000000000000, 0.017452406437, 0.034899496703, 0.052335956243, 0.069756473744, 0.087155742748, 0.104528463268, 0.121869343405, 0.139173100960, 0.156434465040, 0.173648177667, 0.190808995377, 0.207911690818, 0.224951054344, 0.241921895600, 0.258819045103, 0.275637355817, 0.292371704723, 0.309016994375, 0.325568154457, 0.342020143326, 0.358367949545, 0.374606593416, 0.390731128489, 0.406736643076, 0.422618261741, 0.438371146789, 0.453990499740, 0.469471562786, 0.484809620246, 0.500000000000, 0.515038074910, 0.529919264233, 0.544639035015, 0.559192903471, 0.573576436351, 0.587785252292, 0.601815023152, 0.615661475326, 0.629320391050, 0.642787609687, 0.656059028991, 0.669130606359, 0.681998360062, 0.694658370459, 0.707106781187, 0.719339800339, 0.731353701619, 0.743144825477, 0.754709580223, 0.766044443119, 0.777145961457, 0.788010753607, 0.798635510047, 0.809016994375, 0.819152044289, 0.829037572555, 0.838670567945, 0.848048096156, 0.857167300702, 0.866025403784, 0.874619707139, 0.882947592859, 0.891006524188, 0.898794046299, 0.906307787037, 0.913545457643, 0.920504853452, 0.927183854567, 0.933580426497, 0.939692620786, 0.945518575599, 0.951056516295, 0.956304755963, 0.961261695938, 0.965925826289, 0.970295726276, 0.974370064785, 0.978147600734, 0.981627183448, 0.984807753012, 0.987688340595, 0.990268068742, 0.992546151641, 0.994521895368, 0.996194698092, 0.997564050260, 0.998629534755, 0.999390827019, 0.999847695156, 1.000000000000, 0.999847695156, 0.999390827019, 0.998629534755, 0.997564050260, 0.996194698092, 0.994521895368, 0.992546151641, 0.990268068742, 0.987688340595, 0.984807753012, 0.981627183448, 0.978147600734, 0.974370064785, 0.970295726276, 0.965925826289, 0.961261695938, 0.956304755963, 0.951056516295, 0.945518575599, 0.939692620786, 0.933580426497, 0.927183854567, 0.920504853452, 0.913545457643, 0.906307787037, 0.898794046299, 0.891006524188, 0.882947592859, 0.874619707139, 0.866025403784, 0.857167300702, 0.848048096156, 0.838670567945, 0.829037572555, 0.819152044289, 0.809016994375, 0.798635510047, 0.788010753607, 0.777145961457, 0.766044443119, 0.754709580223, 0.743144825477, 0.731353701619, 0.719339800339, 0.707106781187, 0.694658370459, 0.681998360062, 0.669130606359, 0.656059028991, 0.642787609687, 0.629320391050, 0.615661475326, 0.601815023152, 0.587785252292, 0.573576436351, 0.559192903471, 0.544639035015, 0.529919264233, 0.515038074910, 0.500000000000, 0.484809620246, 0.469471562786, 0.453990499740, 0.438371146789, 0.422618261741, 0.406736643076, 0.390731128489, 0.374606593416, 0.358367949545, 0.342020143326, 0.325568154457, 0.309016994375, 0.292371704723, 0.275637355817, 0.258819045103, 0.241921895600, 0.224951054344, 0.207911690818, 0.190808995377, 0.173648177667, 0.156434465040, 0.139173100960, 0.121869343405, 0.104528463268, 0.087155742748, 0.069756473744, 0.052335956243, 0.034899496703, 0.017452406437}

/*
Rotate a fyne.Position around a center point given a specific angle in degrees!
*/
func rotatePosition(centerX, centerY float64, point *fyne.Position, angle int) {
	dx := float64(point.X) - centerX // Cal deltas
	dy := float64(point.Y) - centerY
	point.X = float32(cos(angle)*dx - sin(angle)*dy + centerX) // Rotate using fast sine values
	point.Y = float32(sin(angle)*dx + cos(angle)*dy + centerY)
}

func rotatePoints(centerX, centerY, x, y float64, angle int) (float64, float64) {
	dx := x - centerX // Cal deltas
	dy := y - centerY
	px := cos(angle)*dx - sin(angle)*dy + centerX // Rotate using fast sine values
	py := sin(angle)*dx + cos(angle)*dy + centerY
	return px, py
}

func ScaleMovable(mov Movable, scale float64) {
	sc := mov.GetSizeAndCenter()
	w := sc.Width * scale
	h := sc.Height * scale
	if w < 1 || h < 1 || h > 10000 || w > 10000 {
		return
	}
	mov.SetSize(fyne.Size{Width: float32(w), Height: float32(h)})
}

func SetSpeedAndDirection(mov Movable, speed float64, angle int) {
	dx := speed * cos(angle)
	dy := speed * sin(angle)
	mov.SetSpeed(dx, dy)
}

func SetSpeedAndTarget(mov, target Movable, speed float64) {
	x, y := target.GetCenter()
	SetSpeedAndTargetPosition(mov, speed, x, y)
}

func SetSpeedAndTargetPosition(mov Movable, speed, x, y float64) {
	xp, yp := mov.GetCenter()
	SetSpeedAndDirection(mov, speed, degreesFromCords(xp, yp, x, y))
}

func degreesFromCords(fx, fy, tx, ty float64) int {
	dx := tx - fx
	dy := ty - fy
	if dx < 0 {
		return int(math.Atan(dy/dx)*radToDeg) + 180
	}
	if dy < 0 {
		return int(math.Atan(dy/dx)*radToDeg) + 360
	}
	return int(math.Atan(dy/dx) * radToDeg)
}

func distanceFromCords(fx, fy, tx, ty float64) float64 {
	dx := tx - fx
	dy := ty - fy
	return math.Abs(math.Sqrt(dx*dx + dy*dy))
}

func scalePoint(centerX, centerY float64, point *fyne.Position, scaleX, scaleY float64) {
	dx := float64(point.X) - centerX // Cal deltas
	dy := float64(point.Y) - centerY
	dx = scaleX * dx
	dy = scaleY * dy
	point.X = float32(centerX + dx)
	point.Y = float32(centerY + dy)
}

func sin(deg int) float64 {
	deg = deg % 360 // Wrap value. Can only be 0..359
	if deg < 0 {
		deg = deg + 360 // if negative make positive
	}
	if deg >= 180 {
		return SIN_COS_TABLE[deg%180] * -1 // Wrap to  0..179, acess array value, return negative value
	}
	return SIN_COS_TABLE[deg] // Wrap to  0..179, acess array value, return positive value
}

func cos(deg int) float64 {
	deg = (deg + 90) % 360 // Convert to sin value, wrap value. Can only be 0..359
	if deg < 0 {
		deg = deg + 360 // if negative make positive
	}
	if deg >= 180 {
		return SIN_COS_TABLE[deg%180] * -1 // Wrap to 0..179, acess array value, return negative value
	}
	return SIN_COS_TABLE[deg] // Wrap to  0..179, acess array value, return positive value
}
