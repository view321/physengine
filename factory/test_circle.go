package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTestCircle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry{
	entity := ecs.World.Create(components.Transform, components.CircleCollider, components.Drawable, components.MassComponent, components.Velocity)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}
	components.MassComponent.Get(entry).Mass = 10
	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 100
	crcl.Restitution = 0.8
	return entry
}