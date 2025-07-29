package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type AABB_Data struct {
	Min         Vec2.Vec2
	Max         Vec2.Vec2
	Restitution float64
}

var AABB_Component = donburi.NewComponentType[AABB_Data]()

func AABBvsAABB(a1, a2 *donburi.Entry) bool {
	tr1 := Transform.Get(a1)
	tr2 := Transform.Get(a2)
	AABB1 := donburi.Get[AABB_Data](a1, AABB_Component)
	AABB2 := donburi.Get[AABB_Data](a2, AABB_Component)

	if (tr1.Pos.X+AABB1.Max.X < tr2.Pos.X+AABB2.Min.X) || (tr1.Pos.X+AABB1.Min.X > tr2.Pos.X+AABB2.Max.X) {
		return false
	}
	if (tr1.Pos.Y+AABB1.Max.Y < tr2.Pos.Y+AABB2.Min.Y) || (tr1.Pos.Y+AABB1.Min.Y > tr2.Pos.Y+AABB2.Max.Y) {
		return false
	}
	return true
}
