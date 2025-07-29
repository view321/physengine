package Vec2

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (v *Vec2) Invert() {
	v.X *= -1
	v.Y *= -1
}
func (v *Vec2) Magnitude() float64 {
	return float64(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}
func (v *Vec2) SquareMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}
func (v *Vec2) Normalize() {
	l := v.Magnitude()
	v.X = v.X / l
	v.Y = v.Y / l
}
func (v Vec2) Mult(nm float64) Vec2 {
	return Vec2{v.X * nm, v.Y * nm}
}
func (v *Vec2) MultUpdate(nm float64) {
	v.X *= nm
	v.Y *= nm
}
func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}
func (v1 *Vec2) AddUpdate(v2 Vec2) {
	v1.X += v2.X
	v1.Y += v2.Y
}
func DotProduct(v1, v2 Vec2) float64 {
	return float64(v1.X*v2.X + v1.Y*v2.Y)
}
func Add(v1, v2 Vec2) Vec2 {
	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}
func Distance(v1, v2 Vec2) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}
