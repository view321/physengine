package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTestSquare(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry{
	entity := ecs.World.Create(components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}
	components.MassComponent.Get(entry).Mass = 10
	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{-100, -100}
	box.Max = Vec2.Vec2{50, 50}
	box.Restitution = 0.8
	return entry
}