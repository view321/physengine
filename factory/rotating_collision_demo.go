package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// CreateRotatingCollisionDemo creates objects to demonstrate rotation-aware collisions
func CreateRotatingCollisionDemo(ecs *ecs.ECS) {
	// Create a rotating square that will collide with other objects
	CreateRotatingSquare(ecs, Vec2.Vec2{X: 300, Y: 200}, Vec2.Vec2{X: -50, Y: 0}, 3.0)

	// Create a rotating circle that will collide with the square
	CreateRotatingCircle(ecs, Vec2.Vec2{X: 100, Y: 200}, Vec2.Vec2{X: 100, Y: 0}, -2.0)

	// Create a stationary object that will be hit by rotating objects
	CreateStationaryObject(ecs, Vec2.Vec2{X: 500, Y: 300}, Vec2.Vec2{X: 0, Y: 0})
}

// CreateRotatingSquare creates a square with rotation and collision
func CreateRotatingSquare(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2, angularVel float64) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, angularVel)

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 8
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a square
	mc.Inertia = mc.Mass * 8000 // Simplified inertia calculation
	mc.InverseInertia = 1 / mc.Inertia

	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{X: -60, Y: -60}
	box.Max = Vec2.Vec2{X: 60, Y: 200}

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.7
	mat.StaticFriction = 0.3
	mat.DynamicFriction = 0.2

	return entry
}

// CreateRotatingCircle creates a circle with rotation and collision
func CreateRotatingCircle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2, angularVel float64) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.CircleCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, angularVel)

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 6
	mc.InverseMass = 1 / mc.Mass

	// Calculate inertia for a circle
	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 70
	mc.Inertia = mc.Mass * crcl.Radius * crcl.Radius * 0.5 // I = 0.5 * m * r^2 for a circle
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.8
	mat.StaticFriction = 0.4
	mat.DynamicFriction = 0.3

	return entry
}

// CreateStationaryObject creates a stationary object for collision testing
func CreateStationaryObject(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 0.0) // No rotation

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1, Y: 1}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 20 // Heavy object
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a square
	mc.Inertia = mc.Mass * 12000
	mc.InverseInertia = 1 / mc.Inertia

	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{X: -80, Y: -80}
	box.Max = Vec2.Vec2{X: 80, Y: 80}

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.5
	mat.StaticFriction = 0.6
	mat.DynamicFriction = 0.4

	return entry
}
