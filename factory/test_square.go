package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTestSquare(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry{
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 2.0) // Add some rotation
	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}
	mc := components.MassComponent.Get(entry)
	mc.Mass = 10
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a square (approximate)
	mc.Inertia = mc.Mass * 10000 // Simplified inertia calculation
	mc.InverseInertia = 1 / mc.Inertia
	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{-100, -50}
	box.Max = Vec2.Vec2{50, 50}
	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.8
	return entry
}