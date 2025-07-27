package components

import (
	"math"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type TransformData struct {
	Pos      Vec2.Vec2
	Rot      float64
	Scale    Vec2.Vec2
	Children []*TransformData
	ID       int
}

func (p TransformData) Order() int {
	return -p.ID
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
