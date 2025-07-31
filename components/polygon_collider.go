package components

import (
	"math"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type PolygonColliderData struct {
	Vertices []Vec2.Vec2
}

var PolygonCollider = donburi.NewComponentType[PolygonColliderData]()

func CalculateCentroid(vertices []Vec2.Vec2) Vec2.Vec2 {
	var centroid Vec2.Vec2
	for _, v := range vertices {
		centroid.X += v.X
		centroid.Y += v.Y
	}
	centroid.X /= float64(len(vertices))
	centroid.Y /= float64(len(vertices))
	return centroid
}

func ProjectPolygon(e *donburi.Entry, axis Vec2.Vec2) (float64, float64) {
	polygon_data := donburi.GetValue[PolygonColliderData](e, PolygonCollider)
	tr := Transform.Get(e)
	for _, vertice := range polygon_data.Vertices {
		vertice.X += tr.Pos.X
		vertice.Y += tr.Pos.Y
	}
	min := Vec2.Scalar(axis, polygon_data.Vertices[0])
	max := min
	for i := 1; i < len(polygon_data.Vertices); i++ {
		d := Vec2.Scalar(axis, polygon_data.Vertices[i])
		if d < min {
			min = d
		} else if d > max {
			max = d
		}
	}
	return min, max
}

func ClosestPoint(e *donburi.Entry, x, y float64) (float64, float64) {
	polygon_data := donburi.GetValue[PolygonColliderData](e, PolygonCollider)
	tr := Transform.Get(e)
	for _, vertice := range polygon_data.Vertices {
		vertice.X += tr.Pos.X
		vertice.Y += tr.Pos.Y
	}
	minDistSquared := math.MaxFloat64
	closestX, closestY := 0.0, 0.0
	for i := 0; i < len(polygon_data.Vertices); i++ {
		dx := x - polygon_data.Vertices[i].X
		dy := y - polygon_data.Vertices[i].Y
		squared_dist := dx*dx + dy*dy
		if squared_dist < minDistSquared {
			closestX = polygon_data.Vertices[i].X
			closestY = polygon_data.Vertices[i].Y
		}
	}
	return closestX, closestY
}
