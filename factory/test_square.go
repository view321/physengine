package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTestSquare(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.GravityTag, components.ForceComponent, components.MaterialComponent, components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent)
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
	// Calculate inertia for a square (realistic value for friction effectiveness)
	// For a square: I = m * (width^2 + height^2) / 12
	// With width=150, height=100: I = 10 * (150^2 + 100^2) / 12 = 10 * 32500 / 12 ≈ 27083
	mc.Inertia = 2700 // Realistic inertia for a square
	mc.InverseInertia = 1 / mc.Inertia
	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{X: -100, Y: -50}
	box.Max = Vec2.Vec2{X: 50, Y: 50}
	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.8
	mat.StaticFriction = 1.0  // Very high static friction
	mat.DynamicFriction = 0.8 // Very high dynamic friction
	return entry
}
