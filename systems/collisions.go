package systems

import (
	"math"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateCollisions(e *ecs.ECS) {
	resolver_entry, _ := components.CollisionResolverComponent.First(e.World)
	resolver_comp := components.CollisionResolverComponent.Get(resolver_entry)
	for num1 := 0; num1 < len(resolver_comp.Physobs); num1++ {
		for num2 := num1 + 1; num2 < len(resolver_comp.Physobs); num2++ {
			ResolveCollisions(resolver_comp.Physobs[num1], resolver_comp.Physobs[num2])
		}
	}
}
func ResolveCollisions(c1, c2 *donburi.Entry) {
	if c1.HasComponent(components.CircleCollider) && c2.HasComponent(components.CircleCollider) {
		if components.CirclesCollide(c1, c2) {
			tr1 := components.Transform.Get(c1)
			tr2 := components.Transform.Get(c2)
			crcl1 := components.CircleCollider.Get(c1)
			crcl2 := components.CircleCollider.Get(c2)
			vl1 := components.Velocity.Get(c1)
			vl2 := components.Velocity.Get(c2)
			ms1 := components.MassComponent.Get(c1)
			ms2 := components.MassComponent.Get(c2)
			rel_vel := vl2.Velocity.Add(vl1.Velocity.Mult(-1))
			normal := tr2.Pos.Add(tr1.Pos.Mult(-1))
			velAlongNormal := Vec2.DotProduct(rel_vel, normal)
			if velAlongNormal > 0 {
				return
			}
			e := math.Min(crcl1.Restitution, crcl2.Restitution)
			j := -(1 + e) * velAlongNormal
			impulse := normal.Mult(j)
			vl1.Velocity = vl1.Velocity.Add(impulse.Mult(-1 / ms1.Mass))
			vl2.Velocity = vl2.Velocity.Add(impulse.Mult(1 / ms2.Mass))

			pos_diff := tr1.Pos.Add(tr2.Pos.Mult(-1))
			penetration := crcl1.Radius + crcl2.Radius - pos_diff.Magnitude()
			correction := normal.Mult(math.Max(penetration-0.01, 0.0) * 0.2 / (1/ms1.Mass + 1/ms2.Mass))
			components.SetPos(c1, tr1.Pos.Add(correction.Mult(-1/ms1.Mass)))
			components.SetPos(c2, tr2.Pos.Add(correction.Mult(1/ms2.Mass)))
		}
	}
}
