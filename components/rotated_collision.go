package components

import (
	"math"
	"physengine/helpers"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

// RotatedAABBvsAABB checks collision between two rotated AABBs
func RotatedAABBvsAABB(a1, a2 *donburi.Entry) (bool, Vec2.Vec2, float64) {
	if !a1.HasComponent(Transform) || !a2.HasComponent(Transform) {
		return false, Vec2.Vec2{}, 0
	}
	if !a1.HasComponent(AABB_Component) || !a2.HasComponent(AABB_Component) {
		return false, Vec2.Vec2{}, 0
	}

	tr1 := Transform.Get(a1)
	tr2 := Transform.Get(a2)
	AABB1 := AABB_Component.Get(a1)
	AABB2 := AABB_Component.Get(a2)

	if tr1 == nil || tr2 == nil || AABB1 == nil || AABB2 == nil {
		return false, Vec2.Vec2{}, 0
	}

	// Get the corners of both AABBs in world space
	corners1 := getRotatedAABBCorners(tr1, AABB1)
	corners2 := getRotatedAABBCorners(tr2, AABB2)

	// Use Separating Axis Theorem (SAT) for rotated AABB collision
	return SatCollision(corners1, corners2)
}

// getRotatedAABBCorners returns the 4 corners of a rotated AABB in world space
func getRotatedAABBCorners(tr *TransformData, aabb *AABB_Data) []Vec2.Vec2 {
	// Calculate center and half-extents
	center := tr.Pos
	halfWidth := (aabb.Max.X - aabb.Min.X) / 2
	halfHeight := (aabb.Max.Y - aabb.Min.Y) / 2

	// Define local corners (before rotation)
	localCorners := []Vec2.Vec2{
		{X: -halfWidth, Y: -halfHeight}, // bottom-left
		{X: halfWidth, Y: -halfHeight},  // bottom-right
		{X: halfWidth, Y: halfHeight},   // top-right
		{X: -halfWidth, Y: halfHeight},  // top-left
	}

	// Rotate and translate corners to world space
	worldCorners := make([]Vec2.Vec2, 4)
	for i, corner := range localCorners {
		rotated := RotatePoint(corner, tr.Rot)
		worldCorners[i] = Vec2.Vec2{
			X: center.X + rotated.X,
			Y: center.Y + rotated.Y,
		}
	}

	return worldCorners
}

// SatCollision performs Separating Axis Theorem collision detection
func SatCollision(corners1, corners2 []Vec2.Vec2) (bool, Vec2.Vec2, float64) {
	// Get axes to test (normals of edges)
	axes := getAxes(corners1, corners2)

	minOverlap := math.Inf(1)
	var collisionNormal Vec2.Vec2

	// Test each axis
	for _, axis := range axes {
		proj1 := projectPolygon(corners1, axis)
		proj2 := projectPolygon(corners2, axis)

		// Check for separation
		if proj1.Max < proj2.Min || proj2.Max < proj1.Min {
			return false, Vec2.Vec2{}, 0
		}

		// Calculate overlap
		overlap := math.Min(proj1.Max-proj2.Min, proj2.Max-proj1.Min)
		if overlap < minOverlap {
			minOverlap = overlap
			collisionNormal = axis
		}
	}

	return true, collisionNormal, minOverlap
}

// getAxes returns the axes to test for SAT collision detection
func getAxes(corners1, corners2 []Vec2.Vec2) []Vec2.Vec2 {
	axes := make([]Vec2.Vec2, 0)

	// Add axes from first polygon
	for i := 0; i < len(corners1); i++ {
		next := (i + 1) % len(corners1)
		edge := Vec2.Vec2{
			X: corners1[next].X - corners1[i].X,
			Y: corners1[next].Y - corners1[i].Y,
		}
		normal := Vec2.Vec2{X: -edge.Y, Y: edge.X}
		normal.Normalize()
		axes = append(axes, normal)
	}

	// Add axes from second polygon
	for i := 0; i < len(corners2); i++ {
		next := (i + 1) % len(corners2)
		edge := Vec2.Vec2{
			X: corners2[next].X - corners2[i].X,
			Y: corners2[next].Y - corners2[i].Y,
		}
		normal := Vec2.Vec2{X: -edge.Y, Y: edge.X}
		normal.Normalize()
		axes = append(axes, normal)
	}

	return axes
}

// projectPolygon projects a polygon onto an axis
func projectPolygon(corners []Vec2.Vec2, axis Vec2.Vec2) struct {
	Min, Max float64
} {
	min := Vec2.DotProduct(corners[0], axis)
	max := min

	for _, corner := range corners {
		proj := Vec2.DotProduct(corner, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}

	return struct {
		Min, Max float64
	}{min, max}
}

// RotatedCirclevsAABB checks collision between a circle and a rotated AABB
func RotatedCirclevsAABB(circle, box *donburi.Entry) (bool, Vec2.Vec2, float64) {
	if !circle.HasComponent(Transform) || !box.HasComponent(Transform) {
		return false, Vec2.Vec2{}, 0
	}
	if !circle.HasComponent(CircleCollider) || !box.HasComponent(AABB_Component) {
		return false, Vec2.Vec2{}, 0
	}

	circleTr := Transform.Get(circle)
	boxTr := Transform.Get(box)
	circleComp := CircleCollider.Get(circle)
	boxComp := AABB_Component.Get(box)

	if circleTr == nil || boxTr == nil || circleComp == nil || boxComp == nil {
		return false, Vec2.Vec2{}, 0
	}

	// Transform circle center to box's local space
	circleToBox := Vec2.Vec2{
		X: circleTr.Pos.X - boxTr.Pos.X,
		Y: circleTr.Pos.Y - boxTr.Pos.Y,
	}

	// Rotate circle center to box's local coordinate system
	cos := math.Cos(-boxTr.Rot)
	sin := math.Sin(-boxTr.Rot)
	localCircle := Vec2.Vec2{
		X: circleToBox.X*cos - circleToBox.Y*sin,
		Y: circleToBox.X*sin + circleToBox.Y*cos,
	}

	// Calculate box half-extents
	halfWidth := (boxComp.Max.X - boxComp.Min.X) / 2
	halfHeight := (boxComp.Max.Y - boxComp.Min.Y) / 2

	// Find closest point on box to circle center
	closestX := helpers.Clamp(localCircle.X, -halfWidth, halfWidth)
	closestY := helpers.Clamp(localCircle.Y, -halfHeight, halfHeight)
	closestPoint := Vec2.Vec2{X: closestX, Y: closestY}

	// Calculate distance from circle center to closest point
	distance := Vec2.Distance(localCircle, closestPoint)

	if distance > circleComp.Radius {
		return false, Vec2.Vec2{}, 0
	}

	// Calculate collision normal in world space
	normal := Vec2.Vec2{
		X: localCircle.X - closestPoint.X,
		Y: localCircle.Y - closestPoint.Y,
	}

	// Rotate normal back to world space
	cos = math.Cos(boxTr.Rot)
	sin = math.Sin(boxTr.Rot)
	worldNormal := Vec2.Vec2{
		X: normal.X*cos - normal.Y*sin,
		Y: normal.X*sin + normal.Y*cos,
	}

	if worldNormal.Magnitude() > 0.001 {
		worldNormal.Normalize()
	} else {
		// If circles are exactly on top of each other, use a default normal
		worldNormal = Vec2.Vec2{X: 1, Y: 0}
	}

	penetration := circleComp.Radius - distance
	return true, worldNormal, penetration
}
