package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTestCircle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.CircleCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, -1.5) // Add some rotation in opposite direction
	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}
	mc := components.MassComponent.Get(entry)
	mc.Mass = 10
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a circle
	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 100
	mc.Inertia = mc.Mass * crcl.Radius * crcl.Radius * 0.5 // I = 0.5 * m * r^2 for a circle
	mc.InverseInertia = 1 / mc.Inertia
	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.8
	return entry
}
