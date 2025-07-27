package components

import (
	"math"
	Vec2 "physengine/helpers/vec2"
	"github.com/yohamta/donburi"
)

type CircleColliderData struct{
	Transform *TransformData
	Radius float64
}

var CircleCollider = donburi.NewComponentType[CircleColliderData]()

func PosInsideCollider(circle *CircleColliderData, pos *Vec2.Vec2) bool{
	return math.Pow(pos.X - circle.Transform.Pos.X, 2) + math.Pow(pos.Y - circle.Transform.Pos.Y, 2) < math.Pow(circle.Radius, 2)
}