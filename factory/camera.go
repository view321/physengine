package factory

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateCamera(ecs *ecs.ECS) *donburi.Entry {
	entity := ecs.World.Create(components.Camera, components.Transform)
	entry := ecs.World.Entry(entity)
	components.Camera.SetValue(entry, components.CameraData{ViewportSizeX: 1000, ViewportSizeY: 1000})

	// Initialize camera transform
	tr := components.Transform.Get(entry)
	tr.Pos = Vec2.Vec2{X: 500, Y: 500} // Center of the screen
	tr.Scale = Vec2.Vec2{X: 1, Y: 1}
	tr.Rot = 0
	tr.ID = 0

	return entry
}
