package components

import (
	"fmt"
	"math"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type TransformData struct {
	Pos      Vec2.Vec2
	Rot      float64
	Scale    Vec2.Vec2
	Children []*TransformData
}

var Transform = donburi.NewComponentType[TransformData]()

// GetWorldPosition calculates the world position of a transform
func GetWorldPosition(transform *TransformData) Vec2.Vec2 {
	return transform.Pos
}

// GetWorldRotation calculates the world rotation of a transform
func GetWorldRotation(transform *TransformData) float64 {
	return transform.Rot
}

// GetWorldScale calculates the world scale of a transform
func GetWorldScale(transform *TransformData) Vec2.Vec2 {
	return transform.Scale
}

func SetTransform(entry *donburi.Entry, new_transform TransformData) {
	old_transform := Transform.Get(entry)
	if old_transform == nil {
		fmt.Println("Object does not have transform component")
		return
	}
	// Calculate deltas
	rotation_delta := new_transform.Rot - old_transform.Rot
	scale_delta_x := new_transform.Scale.X / old_transform.Scale.X
	scale_delta_y := new_transform.Scale.Y / old_transform.Scale.Y

	for _, child := range old_transform.Children {
		// Calculate child position relative to parent
		relative_x := child.Pos.X - old_transform.Pos.X
		relative_y := child.Pos.Y - old_transform.Pos.Y

		// Apply scale changes to relative position
		relative_x *= scale_delta_x
		relative_y *= scale_delta_y

		// Apply rotation changes to relative position
		rotated_x := relative_x*math.Cos(rotation_delta) - relative_y*math.Sin(rotation_delta)
		rotated_y := relative_x*math.Sin(rotation_delta) + relative_y*math.Cos(rotation_delta)

		// Set child position to new parent position + rotated relative position
		child.Pos.X = new_transform.Pos.X + rotated_x
		child.Pos.Y = new_transform.Pos.Y + rotated_y

		// Update child rotation and scale
		child.Rot += rotation_delta
		child.Scale.X *= scale_delta_x
		child.Scale.Y *= scale_delta_y
	}

	old_transform.Pos = new_transform.Pos
	old_transform.Scale = new_transform.Scale
	old_transform.Rot = new_transform.Rot
}

func SetPos(entry *donburi.Entry, new_pos Vec2.Vec2) {
	old_transform := Transform.Get(entry)
	for _, child := range old_transform.Children {
		child.Pos.X = child.Pos.X + new_pos.X - old_transform.Pos.X
		child.Pos.Y = child.Pos.Y + new_pos.Y - old_transform.Pos.Y
	}
	old_transform.Pos = new_pos
}

func ChangePos(e *donburi.Entry, pos_diff Vec2.Vec2){
	tr := Transform.Get(e)
	for _, child := range tr.Children{
		child.Pos.AddUpdate(pos_diff)
	}
	tr.Pos.AddUpdate(pos_diff)
}

func SetRot(entry *donburi.Entry, new_rot float64) {
	old_transform := Transform.Get(entry)
	rotation_delta := new_rot - old_transform.Rot

	for _, child := range old_transform.Children {
		// Calculate child position relative to parent
		relative_x := child.Pos.X - old_transform.Pos.X
		relative_y := child.Pos.Y - old_transform.Pos.Y

		// Rotate the relative position around parent's origin
		rotated_x := relative_x*math.Cos(rotation_delta) - relative_y*math.Sin(rotation_delta)
		rotated_y := relative_x*math.Sin(rotation_delta) + relative_y*math.Cos(rotation_delta)

		// Set child position to parent position + rotated relative position
		child.Pos.X = old_transform.Pos.X + rotated_x
		child.Pos.Y = old_transform.Pos.Y + rotated_y

		// Child rotation should inherit parent's rotation change
		child.Rot += rotation_delta
	}

	old_transform.Rot = new_rot
}
func Rotate(entry *donburi.Entry, rot float64){
	old_transform := Transform.Get(entry)
	rotation_delta := rot

	for _, child := range old_transform.Children {
		// Calculate child position relative to parent
		relative_x := child.Pos.X - old_transform.Pos.X
		relative_y := child.Pos.Y - old_transform.Pos.Y

		// Rotate the relative position around parent's origin
		rotated_x := relative_x*math.Cos(rotation_delta) - relative_y*math.Sin(rotation_delta)
		rotated_y := relative_x*math.Sin(rotation_delta) + relative_y*math.Cos(rotation_delta)

		// Set child position to parent position + rotated relative position
		child.Pos.X = old_transform.Pos.X + rotated_x
		child.Pos.Y = old_transform.Pos.Y + rotated_y

		// Child rotation should inherit parent's rotation change
		child.Rot += rotation_delta
	}

	old_transform.Rot = old_transform.Rot + rot
}

// GetRotationMatrix returns a 2x2 rotation matrix for the given angle
func GetRotationMatrix(angle float64) [4]float64 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return [4]float64{cos, -sin, sin, cos}
}

// RotatePoint rotates a point around the origin by the given angle
func RotatePoint(point Vec2.Vec2, angle float64) Vec2.Vec2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vec2.Vec2{
		X: point.X*cos - point.Y*sin,
		Y: point.X*sin + point.Y*cos,
	}
}

// RotatePointAround rotates a point around a center point by the given angle
func RotatePointAround(point Vec2.Vec2, center Vec2.Vec2, angle float64) Vec2.Vec2 {
	// Translate to origin
	translated := Vec2.Vec2{X: point.X - center.X, Y: point.Y - center.Y}
	// Rotate
	rotated := RotatePoint(translated, angle)
	// Translate back
	return Vec2.Vec2{X: rotated.X + center.X, Y: rotated.Y + center.Y}
}