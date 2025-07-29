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
	if !a1.HasComponent(Transform) || !a2.HasComponent(Transform) {
		return false
	}
	if !a1.HasComponent(AABB_Component) || !a2.HasComponent(AABB_Component) {
		return false
	}

	tr1 := Transform.Get(a1)
	tr2 := Transform.Get(a2)
	AABB1 := donburi.Get[AABB_Data](a1, AABB_Component)
	AABB2 := donburi.Get[AABB_Data](a2, AABB_Component)

	if tr1 == nil || tr2 == nil || AABB1 == nil || AABB2 == nil {
		return false
	}

	// Calculate world positions of AABB bounds
	a1_min_x := tr1.Pos.X + AABB1.Min.X
	a1_max_x := tr1.Pos.X + AABB1.Max.X
	a1_min_y := tr1.Pos.Y + AABB1.Min.Y
	a1_max_y := tr1.Pos.Y + AABB1.Max.Y

	a2_min_x := tr2.Pos.X + AABB2.Min.X
	a2_max_x := tr2.Pos.X + AABB2.Max.X
	a2_min_y := tr2.Pos.Y + AABB2.Min.Y
	a2_max_y := tr2.Pos.Y + AABB2.Max.Y

	// Check for overlap on both axes
	if (a1_max_x < a2_min_x) || (a1_min_x > a2_max_x) {
		return false
	}
	if (a1_max_y < a2_min_y) || (a1_min_y > a2_max_y) {
		return false
	}
	return true
}
