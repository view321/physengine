package vector2

import "math"

type Vector2 struct {
	X   float64
	Y   float64
}
func (v *Vector2) init (x, y float64){
	v.X = x
	v.Y = y
}
func (v *Vector2) Invert() {
	v.X *= -1
	v.Y *= -1
}
func (v *Vector2) Magnitude() float64 {
	return float64(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}
func (v *Vector2) SquareMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}
func (v *Vector2) Normalize() {
	l := v.Magnitude()
	v.X = v.X / l
	v.Y = v.Y / l
}
func (v Vector2) Mult(nm float64) Vector2 {
	return Vector2{v.X * nm, v.Y * nm}
}
func (v *Vector2) MultUpdate(nm float64) {
	v.X *= nm
	v.Y *= nm
}
func (v1 Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{v1.X + v2.X, v1.Y + v2.Y}
}
func (v1 *Vector2) AddUpdate(v2 Vector2) {
	v1.X += v2.X
	v1.Y += v2.Y
}
func Scalar(v1, v2 Vector2) float64 {
	return float64(v1.X*v2.X + v1.Y*v2.Y)
}
func Add(v1, v2 Vector2) Vector2{
	return Vector2{v1.X+v2.X, v1.Y+v2.Y}
}