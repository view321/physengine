package systems

import (
	"math"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateImprovedCollisions(e *ecs.ECS) {
	resolver_entry, _ := components.CollisionResolverComponent.First(e.World)
	if resolver_entry == nil {
		return
	}
	resolver_comp := components.CollisionResolverComponent.Get(resolver_entry)

	// Update the physics objects list dynamically
	query := donburi.NewQuery(filter.Or(filter.Contains(components.CircleCollider), filter.Contains(components.AABB_Component), filter.Contains(components.PolygonCollider)))
	resolver_comp.Physobs = nil
	for phys_entry := range query.Iter(e.World) {
		resolver_comp.Physobs = append(resolver_comp.Physobs, phys_entry)
	}

	for num1 := 0; num1 < len(resolver_comp.Physobs); num1++ {
		for num2 := num1 + 1; num2 < len(resolver_comp.Physobs); num2++ {
			ResolveImprovedCollisions(resolver_comp.Physobs[num1], resolver_comp.Physobs[num2])
		}
	}
}

func ResolveImprovedCollisions(e1, e2 *donburi.Entry) {
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

			// Normalize the normal vector with better numerical stability
			if distance > 0.001 {
				normal = normal.Mult(1.0 / distance)
			} else {
				// If circles are exactly on top of each other, use a default normal
				normal = Vec2.Vec2{X: 1, Y: 0}
			}

			// Calculate collision point (weighted by radius)
			totalRadius := crcl1.Radius + crcl2.Radius
			ratio1 := crcl2.Radius / totalRadius
			ratio2 := crcl1.Radius / totalRadius
			collisionPoint := Vec2.Vec2{
				X: tr1.Pos.X*ratio2 + tr2.Pos.X*ratio1,
				Y: tr1.Pos.Y*ratio2 + tr2.Pos.Y*ratio1,
			}

			var j float64 = ResolveWithImprovedAngularImpulse(e1, e2, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			penetration := totalRadius - distance
			ImprovedPositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Rotated AABB vs Rotated AABB
	if e1.HasComponent(components.AABB_Component) && e2.HasComponent(components.AABB_Component) {
		colliding, normal, penetration := components.RotatedAABBvsAABB(e1, e2)
		if colliding {
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision point (center of overlap)
			tr1 := components.Transform.Get(e1)
			tr2 := components.Transform.Get(e2)
			collisionPoint := Vec2.Vec2{
				X: (tr1.Pos.X + tr2.Pos.X) / 2,
				Y: (tr1.Pos.Y + tr2.Pos.Y) / 2,
			}

			var j float64 = ResolveWithImprovedAngularImpulse(e1, e2, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			ImprovedPositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
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

			var j float64 = ResolveWithImprovedAngularImpulse(box, circle, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			ImprovedPositionalCorrection(box, circle, normal, penetration, 0.2)
			ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Polygon vs Polygon
	if e1.HasComponent(components.PolygonCollider) && e2.HasComponent(components.PolygonCollider) {
		colliding, normal, penetration := components.PolygonvsPolygon(e1, e2)
		if colliding {
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision point (center of overlap)
			tr1 := components.Transform.Get(e1)
			tr2 := components.Transform.Get(e2)
			collisionPoint := Vec2.Vec2{
				X: (tr1.Pos.X + tr2.Pos.X) / 2,
				Y: (tr1.Pos.Y + tr2.Pos.Y) / 2,
			}

			var j float64 = ResolveWithImprovedAngularImpulse(e1, e2, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			ImprovedPositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Polygon vs Circle
	if (e1.HasComponent(components.PolygonCollider) && e2.HasComponent(components.CircleCollider)) ||
		(e2.HasComponent(components.PolygonCollider) && e1.HasComponent(components.CircleCollider)) {
		var poly *donburi.Entry
		var circle *donburi.Entry
		if e1.HasComponent(components.CircleCollider) {
			circle = e1
			poly = e2
		} else {
			circle = e2
			poly = e1
		}

		colliding, normal, penetration := components.PolygonvsCircle(poly, circle)
		if colliding {
			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate collision point
			polyTr := components.Transform.Get(poly)
			circleTr := components.Transform.Get(circle)
			collisionPoint := Vec2.Vec2{
				X: (polyTr.Pos.X + circleTr.Pos.X) / 2,
				Y: (polyTr.Pos.Y + circleTr.Pos.Y) / 2,
			}

			var j float64 = ResolveWithImprovedAngularImpulse(poly, circle, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
			ImprovedPositionalCorrection(poly, circle, normal, penetration, 0.2)
			ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
		}
	}

	// Polygon vs AABB (treat AABB as a polygon)
	if (e1.HasComponent(components.PolygonCollider) && e2.HasComponent(components.AABB_Component)) ||
		(e2.HasComponent(components.PolygonCollider) && e1.HasComponent(components.AABB_Component)) {
		var poly *donburi.Entry
		var box *donburi.Entry
		if e1.HasComponent(components.AABB_Component) {
			box = e1
			poly = e2
		} else {
			box = e2
			poly = e1
		}

		// Convert AABB to polygon vertices for collision detection
		boxTr := components.Transform.Get(box)
		boxComp := components.AABB_Component.Get(box)
		polyTr := components.Transform.Get(poly)
		polyComp := components.PolygonCollider.Get(poly)

		if boxTr != nil && boxComp != nil && polyTr != nil && polyComp != nil {
			// Create AABB vertices
			halfWidth := (boxComp.Max.X - boxComp.Min.X) / 2
			halfHeight := (boxComp.Max.Y - boxComp.Min.Y) / 2
			aabbVertices := []Vec2.Vec2{
				{X: -halfWidth, Y: -halfHeight},
				{X: halfWidth, Y: -halfHeight},
				{X: halfWidth, Y: halfHeight},
				{X: -halfWidth, Y: halfHeight},
			}

			// Transform AABB vertices to world space
			worldAABBVertices := make([]Vec2.Vec2, 4)
			for i, vertex := range aabbVertices {
				rotated := components.RotatePoint(vertex, boxTr.Rot)
				worldAABBVertices[i] = Vec2.Vec2{
					X: boxTr.Pos.X + rotated.X,
					Y: boxTr.Pos.Y + rotated.Y,
				}
			}

			// Get polygon vertices in world space
			polyVertices := components.GetWorldVertices(poly)
			if polyVertices != nil {
				colliding, normal, penetration := components.SatCollision(polyVertices, worldAABBVertices)
				if colliding {
					mat1 := components.MaterialComponent.Get(e1)
					mat2 := components.MaterialComponent.Get(e2)

					// Calculate collision point
					collisionPoint := Vec2.Vec2{
						X: (polyTr.Pos.X + boxTr.Pos.X) / 2,
						Y: (polyTr.Pos.Y + boxTr.Pos.Y) / 2,
					}

					var j float64 = ResolveWithImprovedAngularImpulse(poly, box, normal, collisionPoint, mat1.Restitution, mat2.Restitution)
					ImprovedPositionalCorrection(poly, box, normal, penetration, 0.2)
					ResolveImprovedFriction(e1, e2, normal, collisionPoint, j)
				}
			}
		}
	}
}

// ResolveWithImprovedAngularImpulse resolves collision with improved numerical stability
func ResolveWithImprovedAngularImpulse(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, res1, res2 float64) float64 {
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

	// Calculate tangential velocities with improved numerical stability
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

	// Total velocity at collision point
	velAtPoint1 := vel1.Velocity.Add(tangentialVel1)
	velAtPoint2 := vel2.Velocity.Add(tangentialVel2)
	relativeVel := velAtPoint2.Add(velAtPoint1.Mult(-1))

	// Calculate velocity along normal
	velAlongNormal := Vec2.DotProduct(normal, relativeVel)

	// Only resolve if objects are moving towards each other
	if velAlongNormal <= -0.001 { // Small threshold to prevent numerical issues
		e := math.Min(res1, res2)
		j := -(1 + e) * velAlongNormal

		// Calculate impulse denominator including angular terms
		denominator := m1.InverseMass + m2.InverseMass

		// Add angular terms to denominator with improved numerical stability
		if angVel1 != nil && m1.InverseInertia > 0 {
			cross1 := r1.X*normal.Y - r1.Y*normal.X
			denominator += cross1 * cross1 * m1.InverseInertia
		}
		if angVel2 != nil && m2.InverseInertia > 0 {
			cross2 := r2.X*normal.Y - r2.Y*normal.X
			denominator += cross2 * cross2 * m2.InverseInertia
		}

		// Prevent division by zero and clamp impulse
		if denominator > 0.001 {
			j /= denominator

			// Clamp impulse to prevent extreme values
			maxImpulse := 1000.0 // Adjust based on your physics scale
			if math.Abs(j) > maxImpulse {
				if j > 0 {
					j = maxImpulse
				} else {
					j = -maxImpulse
				}
			}

			// Apply linear impulse
			impulse := normal.Mult(j)
			vel1.Velocity = vel1.Velocity.Add(impulse.Mult(-m1.InverseMass))
			vel2.Velocity = vel2.Velocity.Add(impulse.Mult(m2.InverseMass))

			// Apply angular impulse with improved stability
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
	}

	return 0
}

// ResolveImprovedFriction resolves friction with improved stability
func ResolveImprovedFriction(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, j float64) {
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

	// Calculate tangent vector with improved stability
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

	if denominator > 0.001 {
		jt /= denominator

		// Apply friction limits with improved stability
		mu := math.Sqrt(mat1.StaticFriction*mat1.StaticFriction + mat2.StaticFriction*mat2.StaticFriction)
		var frictionImpulse Vec2.Vec2

		if math.Abs(jt) < j*mu {
			frictionImpulse = tangent.Mult(jt)
		} else {
			dynamicFriction := math.Sqrt(mat1.DynamicFriction*mat1.DynamicFriction + mat2.DynamicFriction*mat2.DynamicFriction)
			frictionImpulse = tangent.Mult(-j * dynamicFriction)
		}

		// Clamp friction impulse
		maxFrictionImpulse := 500.0 // Adjust based on your physics scale
		if frictionImpulse.Magnitude() > maxFrictionImpulse {
			frictionImpulse.Normalize()
			frictionImpulse = frictionImpulse.Mult(maxFrictionImpulse)
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
}

// ImprovedPositionalCorrection prevents objects from pulling towards each other
func ImprovedPositionalCorrection(e1, e2 *donburi.Entry, n Vec2.Vec2, penetration_depth, percent float64) {
	if penetration_depth < 0.001 { // Small threshold to prevent unnecessary corrections
		return
	}
	m1 := components.MassComponent.Get(e1)
	m2 := components.MassComponent.Get(e2)

	// Safety checks to prevent null pointer issues
	if m1 == nil || m2 == nil {
		return
	}

	// Calculate correction with improved stability
	totalInverseMass := m1.InverseMass + m2.InverseMass
	if totalInverseMass < 0.001 {
		return
	}

	correction := n.Mult(percent * penetration_depth / totalInverseMass)

	// Clamp correction to prevent extreme values
	maxCorrection := 10.0 // Adjust based on your physics scale
	if correction.Magnitude() > maxCorrection {
		correction.Normalize()
		correction = correction.Mult(maxCorrection)
	}

	components.ChangePos(e1, correction.Mult(-m1.InverseMass))
	components.ChangePos(e2, correction.Mult(m2.InverseMass))
}
