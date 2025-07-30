package systems

import (
	"math"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateRotatedCollisions(e *ecs.ECS) {
	resolver_entry, _ := components.CollisionResolverComponent.First(e.World)
	if resolver_entry == nil {
		return
	}
	resolver_comp := components.CollisionResolverComponent.Get(resolver_entry)

	// Update the physics objects list dynamically
	query := donburi.NewQuery(filter.Or(filter.Contains(components.CircleCollider), filter.Contains(components.AABB_Component)))
	resolver_comp.Physobs = nil
	for phys_entry := range query.Iter(e.World) {
		resolver_comp.Physobs = append(resolver_comp.Physobs, phys_entry)
	}

	for num1 := 0; num1 < len(resolver_comp.Physobs); num1++ {
		for num2 := num1 + 1; num2 < len(resolver_comp.Physobs); num2++ {
			ResolveRotatedCollisions(resolver_comp.Physobs[num1], resolver_comp.Physobs[num2])
		}
	}
}

func ResolveRotatedCollisions(e1, e2 *donburi.Entry) {
	// Circle vs Circle (rotation doesn't affect circle collision)
	if e1.HasComponent(components.CircleCollider) && e2.HasComponent(components.CircleCollider) {
		if components.CirclesCollide(e1, e2) {
			tr1 := components.Transform.Get(e1)
			tr2 := components.Transform.Get(e2)
			crcl1 := components.CircleCollider.Get(e1)
			crcl2 := components.CircleCollider.Get(e2)
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision normal (from e1 to e2)
			normal := tr2.Pos.Add(tr1.Pos.Mult(-1))
			distance := normal.Magnitude()

			// Normalize the normal vector
			if distance > 0.001 {
				normal = normal.Mult(1.0 / distance)
			} else {
				normal = Vec2.Vec2{X: 1, Y: 0}
			}

			// Calculate collision point (midpoint between centers)
			collisionPoint := Vec2.Vec2{
				X: (tr1.Pos.X + tr2.Pos.X) / 2,
				Y: (tr1.Pos.Y + tr2.Pos.Y) / 2,
			}

			var j float64 = ResolveWithAngularImpulse(e1, e2, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			penetration := (crcl1.Radius + crcl2.Radius) - distance
			PositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveRotatedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Rotated AABB vs Rotated AABB
	if e1.HasComponent(components.AABB_Component) && e2.HasComponent(components.AABB_Component) {
		colliding, normal, penetration := components.RotatedAABBvsAABB(e1, e2)
		if colliding {
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision point (approximate as center of overlap)
			tr1 := components.Transform.Get(e1)
			tr2 := components.Transform.Get(e2)
			collisionPoint := Vec2.Vec2{
				X: (tr1.Pos.X + tr2.Pos.X) / 2,
				Y: (tr1.Pos.Y + tr2.Pos.Y) / 2,
			}

			var j float64 = ResolveWithAngularImpulse(e1, e2, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			PositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveRotatedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Circle vs Rotated AABB
	if (e1.HasComponent(components.AABB_Component) && e2.HasComponent(components.CircleCollider)) ||
		(e2.HasComponent(components.AABB_Component) && e1.HasComponent(components.CircleCollider)) {
		var box *donburi.Entry
		var circle *donburi.Entry
		if e1.HasComponent(components.CircleCollider) {
			circle = e1
			box = e2
		} else {
			circle = e2
			box = e1
		}

		colliding, normal, penetration := components.RotatedCirclevsAABB(circle, box)
		if colliding {
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision point
			circleTr := components.Transform.Get(circle)
			boxTr := components.Transform.Get(box)
			collisionPoint := Vec2.Vec2{
				X: (circleTr.Pos.X + boxTr.Pos.X) / 2,
				Y: (circleTr.Pos.Y + boxTr.Pos.Y) / 2,
			}

			var j float64 = ResolveWithAngularImpulse(box, circle, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			PositionalCorrection(box, circle, normal, penetration, 0.2)
			ResolveRotatedFriction(e1, e2, normal, collisionPoint, j)
		}
	}
}

// ResolveWithAngularImpulse resolves collision with both linear and angular impulse
func ResolveWithAngularImpulse(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, res1, res2 float64) float64 {
	vel1 := components.Velocity.Get(e1)
	vel2 := components.Velocity.Get(e2)
	angVel1 := components.AngularVelocity.Get(e1)
	angVel2 := components.AngularVelocity.Get(e2)
	m1 := components.MassComponent.Get(e1)
	m2 := components.MassComponent.Get(e2)

	if vel1 == nil || vel2 == nil || m1 == nil || m2 == nil {
		return 0
	}

	// Calculate relative velocity at collision point
	tr1 := components.Transform.Get(e1)
	tr2 := components.Transform.Get(e2)

	// Vector from center to collision point
	r1 := Vec2.Vec2{X: collisionPoint.X - tr1.Pos.X, Y: collisionPoint.Y - tr1.Pos.Y}
	r2 := Vec2.Vec2{X: collisionPoint.X - tr2.Pos.X, Y: collisionPoint.Y - tr2.Pos.Y}

	// Calculate tangential velocities
	var tangentialVel1, tangentialVel2 Vec2.Vec2
	if angVel1 != nil {
		// v_tangential = ω × r
		tangentialVel1 = Vec2.Vec2{
			X: -angVel1.AngularVelocity * r1.Y,
			Y: angVel1.AngularVelocity * r1.X,
		}
	}
	if angVel2 != nil {
		tangentialVel2 = Vec2.Vec2{
			X: -angVel2.AngularVelocity * r2.Y,
			Y: angVel2.AngularVelocity * r2.X,
		}
	}

	// Total velocity at collision point
	velAtPoint1 := vel1.Velocity.Add(tangentialVel1)
	velAtPoint2 := vel2.Velocity.Add(tangentialVel2)
	relativeVel := velAtPoint2.Add(velAtPoint1.Mult(-1))

	// Calculate velocity along normal
	velAlongNormal := Vec2.DotProduct(normal, relativeVel)

	if velAlongNormal <= 0 {
		e := math.Min(res1, res2)
		j := -(1 + e) * velAlongNormal

		// Calculate impulse denominator including angular terms
		denominator := m1.InverseMass + m2.InverseMass

		// Add angular terms to denominator
		if angVel1 != nil && m1.InverseInertia > 0 {
			// (r1 × normal)² / I1
			cross1 := r1.X*normal.Y - r1.Y*normal.X
			denominator += cross1 * cross1 * m1.InverseInertia
		}
		if angVel2 != nil && m2.InverseInertia > 0 {
			// (r2 × normal)² / I2
			cross2 := r2.X*normal.Y - r2.Y*normal.X
			denominator += cross2 * cross2 * m2.InverseInertia
		}

		j /= denominator

		// Apply linear impulse
		impulse := normal.Mult(j)
		vel1.Velocity = vel1.Velocity.Add(impulse.Mult(-m1.InverseMass))
		vel2.Velocity = vel2.Velocity.Add(impulse.Mult(m2.InverseMass))

		// Apply angular impulse
		if angVel1 != nil && m1.InverseInertia > 0 {
			cross1 := r1.X*normal.Y - r1.Y*normal.X
			angVel1.AngularVelocity -= j * cross1 * m1.InverseInertia
		}
		if angVel2 != nil && m2.InverseInertia > 0 {
			cross2 := r2.X*normal.Y - r2.Y*normal.X
			angVel2.AngularVelocity += j * cross2 * m2.InverseInertia
		}

		return j
	}

	return 0
}

// ResolveRotatedFriction resolves friction with angular effects
func ResolveRotatedFriction(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, j float64) {
	vel1 := components.Velocity.Get(e1)
	vel2 := components.Velocity.Get(e2)
	angVel1 := components.AngularVelocity.Get(e1)
	angVel2 := components.AngularVelocity.Get(e2)
	m1 := components.MassComponent.Get(e1)
	m2 := components.MassComponent.Get(e2)
	mat1 := components.MaterialComponent.Get(e1)
	mat2 := components.MaterialComponent.Get(e2)

	if vel1 == nil || vel2 == nil || m1 == nil || m2 == nil {
		return
	}

	tr1 := components.Transform.Get(e1)
	tr2 := components.Transform.Get(e2)

	// Vector from center to collision point
	r1 := Vec2.Vec2{X: collisionPoint.X - tr1.Pos.X, Y: collisionPoint.Y - tr1.Pos.Y}
	r2 := Vec2.Vec2{X: collisionPoint.X - tr2.Pos.X, Y: collisionPoint.Y - tr2.Pos.Y}

	// Calculate relative velocity at collision point
	var tangentialVel1, tangentialVel2 Vec2.Vec2
	if angVel1 != nil {
		tangentialVel1 = Vec2.Vec2{
			X: -angVel1.AngularVelocity * r1.Y,
			Y: angVel1.AngularVelocity * r1.X,
		}
	}
	if angVel2 != nil {
		tangentialVel2 = Vec2.Vec2{
			X: -angVel2.AngularVelocity * r2.Y,
			Y: angVel2.AngularVelocity * r2.X,
		}
	}

	velAtPoint1 := vel1.Velocity.Add(tangentialVel1)
	velAtPoint2 := vel2.Velocity.Add(tangentialVel2)
	relativeVel := velAtPoint2.Add(velAtPoint1.Mult(-1))

	// Calculate tangent vector
	tangent := relativeVel.Add(normal.Mult(-1 * Vec2.DotProduct(relativeVel, normal)))
	if tangent.Magnitude() < 0.001 {
		return
	}
	tangent.Normalize()

	// Calculate tangential impulse
	jt := -Vec2.DotProduct(relativeVel, tangent)

	// Calculate impulse denominator including angular terms
	denominator := m1.InverseMass + m2.InverseMass

	if angVel1 != nil && m1.InverseInertia > 0 {
		cross1 := r1.X*tangent.Y - r1.Y*tangent.X
		denominator += cross1 * cross1 * m1.InverseInertia
	}
	if angVel2 != nil && m2.InverseInertia > 0 {
		cross2 := r2.X*tangent.Y - r2.Y*tangent.X
		denominator += cross2 * cross2 * m2.InverseInertia
	}

	jt /= denominator

	// Apply friction limits
	mu := math.Sqrt(mat1.StaticFriction*mat1.StaticFriction + mat2.StaticFriction*mat2.StaticFriction)
	var frictionImpulse Vec2.Vec2

	if math.Abs(jt) < j*mu {
		frictionImpulse = tangent.Mult(jt)
	} else {
		dynamicFriction := math.Sqrt(mat1.DynamicFriction*mat1.DynamicFriction + mat2.DynamicFriction*mat2.DynamicFriction)
		frictionImpulse = tangent.Mult(-j * dynamicFriction)
	}

	// Apply linear friction impulse
	vel1.Velocity.AddUpdate(frictionImpulse.Mult(-m1.InverseMass))
	vel2.Velocity.AddUpdate(frictionImpulse.Mult(m2.InverseMass))

	// Apply angular friction impulse
	if angVel1 != nil && m1.InverseInertia > 0 {
		cross1 := r1.X*frictionImpulse.Y - r1.Y*frictionImpulse.X
		angVel1.AngularVelocity -= cross1 * m1.InverseInertia
	}
	if angVel2 != nil && m2.InverseInertia > 0 {
		cross2 := r2.X*frictionImpulse.Y - r2.Y*frictionImpulse.X
		angVel2.AngularVelocity += cross2 * m2.InverseInertia
	}
}
