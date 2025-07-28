package factory

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateCamera(ecs *ecs.ECS) *donburi.Entry {
	entity := ecs.World.Create(components.Camera, components.Transform)
	entry := ecs.World.Entry(entity)
	components.Camera.SetValue(entry, components.CameraData{
		ViewportSizeX:      1000,
		ViewportSizeY:      1000,
		LastScreenMousePos: Vec2.Vec2{X: 500, Y: 500}, // Initialize to center of screen
		Op:                 ebiten.DrawImageOptions{},
		Zoom:               Vec2.Vec2{X: 0.5, Y: 0.5},
	})

	// Initialize camera transform
	tr := components.Transform.Get(entry)
	tr.Pos = Vec2.Vec2{X: 0, Y: 0} // Center of the screen
	tr.Scale = Vec2.Vec2{X: 1, Y: 1}
	tr.Rot = 0

	return entry
}
