package systems

import (
	"math"
	"physengine/components"
	"physengine/helpers"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateCollisions(e *ecs.ECS) {
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
			ResolveCollisions(resolver_comp.Physobs[num1], resolver_comp.Physobs[num2])
		}
	}
}

func ResolveCollisions(e1, e2 *donburi.Entry) {
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

			// Normalize the normal vector to prevent extreme velocity changes
			if distance > 0.001 {
				normal = normal.Mult(1.0 / distance)
			} else {
				// If circles are exactly on top of each other, use a default normal
				normal = Vec2.Vec2{X: 1, Y: 0}
			}

			var j float64 = ResolveWithData(e1, e2, normal, mat1.Restitution, mat2.Restitution)
			penetration := (crcl1.Radius + crcl2.Radius) - distance
			PositionalCorrection(e1, e2, normal, penetration, 0.2)
			ResolveFriction(e1, e2, normal, j)
		}
	}
	if e1.HasComponent(components.AABB_Component) && e2.HasComponent(components.AABB_Component) {
		if components.AABBvsAABB(e1, e2) {
			tr1 := components.Transform.Get(e1)
			tr2 := components.Transform.Get(e2)

			// Calculate the vector from e1 to e2
			delta := tr2.Pos.Add(tr1.Pos.Mult(-1))

			var normal Vec2.Vec2
			var penetration float64

			aabb1 := components.AABB_Component.Get(e1)
			aabb2 := components.AABB_Component.Get(e2)

			mat1 := components.MaterialComponent.Get(e1)
			mat2 := components.MaterialComponent.Get(e2)

			// Calculate half-widths and half-heights
			a_width := (aabb1.Max.X - aabb1.Min.X) / 2
			a_height := (aabb1.Max.Y - aabb1.Min.Y) / 2
			b_width := (aabb2.Max.X - aabb2.Min.X) / 2
			b_height := (aabb2.Max.Y - aabb2.Min.Y) / 2

			// Calculate overlap on each axis
			x_overlap := a_width + b_width - math.Abs(delta.X)
			y_overlap := a_height + b_height - math.Abs(delta.Y)

			// Check if there's a collision
			if x_overlap > 0 && y_overlap > 0 {
				// Determine which axis has the smallest overlap (this is the collision normal)
				if x_overlap < y_overlap {
					// X-axis collision
					if delta.X > 0 {
						// e2 is to the right of e1, normal should point right (from e1 to e2)
						normal = Vec2.Vec2{X: 1, Y: 0}
					} else {
						// e2 is to the left of e1, normal should point left (from e1 to e2)
						normal = Vec2.Vec2{X: -1, Y: 0}
					}
					penetration = x_overlap
				} else {
					// Y-axis collision
					if delta.Y > 0 {
						// e2 is above e1, normal should point up (from e1 to e2)
						normal = Vec2.Vec2{X: 0, Y: 1}
					} else {
						// e2 is below e1, normal should point down (from e1 to e2)
						normal = Vec2.Vec2{X: 0, Y: -1}
					}
					penetration = y_overlap
				}

				var j float64 = ResolveWithData(e1, e2, normal, mat1.Restitution, mat2.Restitution)
				PositionalCorrection(e1, e2, normal, penetration, 0.2)
				ResolveFriction(e1, e2, normal, j)
			}
		}
	}
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

		box_col := components.AABB_Component.Get(box)
		box_tr := components.Transform.Get(box)
		circle_tr := components.Transform.Get(circle)
		circle_comp := components.CircleCollider.Get(circle)

		mat1 := components.MaterialComponent.Get(e1)
		mat2 := components.MaterialComponent.Get(e2)
		// Vector from box center to circle center
		circle_to_box := circle_tr.Pos.Add(box_tr.Pos.Mult(-1))

		// Find the closest point on the box to the circle center
		x_extent := (box_col.Max.X - box_col.Min.X) / 2
		y_extent := (box_col.Max.Y - box_col.Min.Y) / 2

		closest_x := helpers.Clamp(circle_to_box.X, -x_extent, x_extent)
		closest_y := helpers.Clamp(circle_to_box.Y, -y_extent, y_extent)
		closest_point := Vec2.Vec2{X: closest_x, Y: closest_y}

		// Vector from closest point to circle center
		closest_to_circle := circle_to_box.Add(closest_point.Mult(-1))
		distance_squared := closest_to_circle.SquareMagnitude()
		radius_squared := circle_comp.Radius * circle_comp.Radius

		// Check if circle overlaps with box
		if distance_squared <= radius_squared {
			distance := closest_to_circle.Magnitude()

			var normal Vec2.Vec2
			var penetration float64

			// Check if circle is inside the box
			if distance < 0.001 {
				// Circle is exactly at the closest point, find the closest edge
				if math.Abs(circle_to_box.X) > math.Abs(circle_to_box.Y) {
					if circle_to_box.X > 0 {
						normal = Vec2.Vec2{X: 1, Y: 0}
					} else {
						normal = Vec2.Vec2{X: -1, Y: 0}
					}
				} else {
					if circle_to_box.Y > 0 {
						normal = Vec2.Vec2{X: 0, Y: 1}
					} else {
						normal = Vec2.Vec2{X: 0, Y: -1}
					}
				}
				penetration = circle_comp.Radius
			} else {
				// Circle is outside or partially inside, use the normal from closest point to circle
				normal = closest_to_circle.Mult(1.0 / distance)
				penetration = circle_comp.Radius - distance
			}

			var j float64 = ResolveWithData(box, circle, normal, mat1.Restitution, mat2.Restitution)
			PositionalCorrection(box, circle, normal, penetration, 0.2)
			ResolveFriction(e1, e2, normal, j)
		}
	}
}
func ResolveWithData(e1, e2 *donburi.Entry, normal Vec2.Vec2, res1, res2 float64) float64 {
	vel1 := donburi.Get[components.VelocityData](e1, components.Velocity)
	vel2 := donburi.Get[components.VelocityData](e2, components.Velocity)
	m1 := donburi.Get[components.MassData](e1, components.MassComponent)
	m2 := donburi.Get[components.MassData](e2, components.MassComponent)

	vel_diff := vel2.Velocity.Add(vel1.Velocity.Mult(-1))
	vel_along_normal := Vec2.DotProduct(normal, vel_diff)
	var j float64
	if vel_along_normal <= 0 {
		e := math.Min(res1, res2)
		j = -(1 + e) * vel_along_normal
		j /= m1.InverseMass + m2.InverseMass
		impulse := normal.Mult(j)
		vel1.Velocity = vel1.Velocity.Add(impulse.Mult(-m1.InverseMass))
		vel2.Velocity = vel2.Velocity.Add(impulse.Mult(m2.InverseMass))
	}
	return j
}
func ResolveFriction(e1, e2 *donburi.Entry, normal Vec2.Vec2, j float64) {
	vel1 := components.Velocity.Get(e1)
	vel2 := components.Velocity.Get(e2)
	m1 := components.MassComponent.Get(e1)
	m2 := components.MassComponent.Get(e2)
	mat1 := components.MaterialComponent.Get(e1)
	mat2 := components.MaterialComponent.Get(e2)
	vel_diff := vel2.Velocity.Add(vel1.Velocity.Mult(-1))
	tangent := vel_diff.Add(normal.Mult(-1 * Vec2.DotProduct(vel_diff, normal)))
	tangent.Normalize()

	jt := -Vec2.DotProduct(vel_diff, tangent)
	jt = jt / (m1.InverseMass + m2.InverseMass)

	mu := math.Sqrt(math.Pow(mat1.StaticFriction, 2) + math.Pow(mat2.StaticFriction, 2))
	var frictionImpulse Vec2.Vec2
	if math.Abs(jt) < j*mu {
		frictionImpulse = tangent.Mult(jt)
	} else {
		var dynamicFriction float64 = math.Sqrt(math.Pow(mat1.DynamicFriction, 2) + math.Pow(mat2.DynamicFriction, 2))
		frictionImpulse = tangent.Mult(-j * dynamicFriction)
	}
	vel1.Velocity.AddUpdate(frictionImpulse.Mult(-m1.InverseMass))
	vel2.Velocity.AddUpdate(frictionImpulse.Mult(m2.InverseMass))
}
func PositionalCorrection(e1, e2 *donburi.Entry, n Vec2.Vec2, penetration_depth, percent float64) {
	if penetration_depth < 0 {
		return
	}
	m1 := components.MassComponent.Get(e1)
	m2 := components.MassComponent.Get(e2)

	// Safety checks to prevent null pointer issues
	if m1 == nil || m2 == nil {
		return
	}

	correction := n.Mult(percent * penetration_depth / (m1.InverseMass + m2.InverseMass))
	components.ChangePos(e1, correction.Mult(-m1.InverseMass))
	components.ChangePos(e2, correction.Mult(m2.InverseMass))
}
