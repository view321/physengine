package components

import (
	"math"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type PolygonColliderData struct {
	Vertices []Vec2.Vec2 // Local vertices relative to center
}

var PolygonCollider = donburi.NewComponentType[PolygonColliderData]()

// GetWorldVertices returns the vertices transformed to world space
func GetWorldVertices(entry *donburi.Entry) []Vec2.Vec2 {
	if !entry.HasComponent(Transform) || !entry.HasComponent(PolygonCollider) {
		return nil
	}

	tr := Transform.Get(entry)
	poly := PolygonCollider.Get(entry)

	if tr == nil || poly == nil {
		return nil
	}

	worldVertices := make([]Vec2.Vec2, len(poly.Vertices))
	for i, vertex := range poly.Vertices {
		// Rotate vertex
		rotated := RotatePoint(vertex, tr.Rot)
		// Translate to world position
		worldVertices[i] = Vec2.Vec2{
			X: tr.Pos.X + rotated.X,
			Y: tr.Pos.Y + rotated.Y,
		}
	}

	return worldVertices
}

// PolygonvsPolygon checks collision between two polygons using SAT
func PolygonvsPolygon(poly1, poly2 *donburi.Entry) (bool, Vec2.Vec2, float64) {
	if !poly1.HasComponent(Transform) || !poly2.HasComponent(Transform) {
		return false, Vec2.Vec2{}, 0
	}
	if !poly1.HasComponent(PolygonCollider) || !poly2.HasComponent(PolygonCollider) {
		return false, Vec2.Vec2{}, 0
	}

	vertices1 := GetWorldVertices(poly1)
	vertices2 := GetWorldVertices(poly2)

	if vertices1 == nil || vertices2 == nil {
		return false, Vec2.Vec2{}, 0
	}

	return SatCollision(vertices1, vertices2)
}

// PolygonvsCircle checks collision between a polygon and a circle
func PolygonvsCircle(poly, circle *donburi.Entry) (bool, Vec2.Vec2, float64) {
	if !poly.HasComponent(Transform) || !circle.HasComponent(Transform) {
		return false, Vec2.Vec2{}, 0
	}
	if !poly.HasComponent(PolygonCollider) || !circle.HasComponent(CircleCollider) {
		return false, Vec2.Vec2{}, 0
	}

	polyTr := Transform.Get(poly)
	circleTr := Transform.Get(circle)
	polyComp := PolygonCollider.Get(poly)
	circleComp := CircleCollider.Get(circle)

	if polyTr == nil || circleTr == nil || polyComp == nil || circleComp == nil {
		return false, Vec2.Vec2{}, 0
	}

	vertices := GetWorldVertices(poly)
	if vertices == nil {
		return false, Vec2.Vec2{}, 0
	}

	// Find the closest point on the polygon to the circle center
	closestPoint := findClosestPointOnPolygon(vertices, circleTr.Pos)
	distance := Vec2.Distance(circleTr.Pos, closestPoint)

	if distance > circleComp.Radius {
		return false, Vec2.Vec2{}, 0
	}

	// Calculate collision normal
	normal := Vec2.Vec2{
		X: circleTr.Pos.X - closestPoint.X,
		Y: circleTr.Pos.Y - closestPoint.Y,
	}

	if normal.Magnitude() > 0.001 {
		normal.Normalize()
	} else {
		// If circle center is exactly on polygon, use a default normal
		normal = Vec2.Vec2{X: 1, Y: 0}
	}

	// Ensure normal points from polygon to circle
	polyToCircle := Vec2.Vec2{
		X: circleTr.Pos.X - polyTr.Pos.X,
		Y: circleTr.Pos.Y - polyTr.Pos.Y,
	}

	// If normal points in wrong direction, flip it
	if Vec2.DotProduct(normal, polyToCircle) < 0 {
		normal = normal.Mult(-1)
	}

	penetration := circleComp.Radius - distance
	return true, normal, penetration
}

// findClosestPointOnPolygon finds the closest point on a polygon to a given point
func findClosestPointOnPolygon(vertices []Vec2.Vec2, point Vec2.Vec2) Vec2.Vec2 {
	if len(vertices) == 0 {
		return point
	}

	closestPoint := vertices[0]
	minDistance := Vec2.Distance(point, closestPoint)

	// Check each edge of the polygon
	for i := 0; i < len(vertices); i++ {
		next := (i + 1) % len(vertices)
		edgeStart := vertices[i]
		edgeEnd := vertices[next]

		// Find closest point on this edge
		closestOnEdge := findClosestPointOnLineSegment(edgeStart, edgeEnd, point)
		distance := Vec2.Distance(point, closestOnEdge)

		if distance < minDistance {
			minDistance = distance
			closestPoint = closestOnEdge
		}
	}

	return closestPoint
}

// findClosestPointOnLineSegment finds the closest point on a line segment to a given point
func findClosestPointOnLineSegment(lineStart, lineEnd, point Vec2.Vec2) Vec2.Vec2 {
	lineVec := Vec2.Vec2{
		X: lineEnd.X - lineStart.X,
		Y: lineEnd.Y - lineStart.Y,
	}

	pointVec := Vec2.Vec2{
		X: point.X - lineStart.X,
		Y: point.Y - lineStart.Y,
	}

	lineLengthSq := lineVec.X*lineVec.X + lineVec.Y*lineVec.Y
	if lineLengthSq == 0 {
		return lineStart
	}

	// Project point onto line
	t := Vec2.DotProduct(pointVec, lineVec) / lineLengthSq

	// Clamp to line segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	return Vec2.Vec2{
		X: lineStart.X + t*lineVec.X,
		Y: lineStart.Y + t*lineVec.Y,
	}
}

// CreateRegularPolygon creates a regular polygon with the given number of sides and radius
func CreateRegularPolygon(sides int, radius float64) []Vec2.Vec2 {
	if sides < 3 {
		return nil
	}

	vertices := make([]Vec2.Vec2, sides)
	angleStep := 2 * math.Pi / float64(sides)

	for i := 0; i < sides; i++ {
		angle := float64(i) * angleStep
		vertices[i] = Vec2.Vec2{
			X: radius * math.Cos(angle),
			Y: radius * math.Sin(angle),
		}
	}

	return vertices
}

// CreateRectangle creates a rectangle with the given width and height
func CreateRectangle(width, height float64) []Vec2.Vec2 {
	halfWidth := width / 2
	halfHeight := height / 2

	return []Vec2.Vec2{
		{X: -halfWidth, Y: -halfHeight}, // bottom-left
		{X: halfWidth, Y: -halfHeight},  // bottom-right
		{X: halfWidth, Y: halfHeight},   // top-right
		{X: -halfWidth, Y: halfHeight},  // top-left
	}
}
