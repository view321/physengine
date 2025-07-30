package components

import (
	"fmt"
	"math"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type CircleColliderData struct {
	Radius            float64
}

var CircleCollider = donburi.NewComponentType[CircleColliderData]()

func PosInsideCollider(entry *donburi.Entry, pos Vec2.Vec2) bool {
	if !entry.HasComponent(Transform) {
		fmt.Println("PosInsideCollider: missing transform component")
		return false
	}
	if !entry.HasComponent(CircleCollider) {
		fmt.Println("PosInsideCollider: missing circle collider component")
		return false
	}
	tr := Transform.Get(entry)
	cc := CircleCollider.Get(entry)

	if tr == nil || cc == nil {
		return false
	}

	return math.Pow(pos.X-tr.Pos.X, 2)+math.Pow(pos.Y-tr.Pos.Y, 2) < math.Pow(cc.Radius, 2)
}

func CirclesCollide(e1, e2 *donburi.Entry) bool {
	if (!e1.HasComponent(Transform)) || (!e2.HasComponent(Transform)) {
		fmt.Println("CirclesCollide: missing transform component")
		return false
	}
	if (!e1.HasComponent(CircleCollider)) || (!e2.HasComponent(CircleCollider)) {
		fmt.Println("CirclesCollide: missing circle collider component")
		return false
	}

	tr1 := Transform.Get(e1)
	tr2 := Transform.Get(e2)
	c1 := CircleCollider.Get(e1)
	c2 := CircleCollider.Get(e2)

	if tr1 == nil || tr2 == nil || c1 == nil || c2 == nil {
		return false
	}

	// Calculate distance between centers
	dx := tr1.Pos.X - tr2.Pos.X
	dy := tr1.Pos.Y - tr2.Pos.Y
	dist_squared := dx*dx + dy*dy
	radius_sum_squared := (c1.Radius + c2.Radius) * (c1.Radius + c2.Radius)

	return dist_squared < radius_sum_squared
}
