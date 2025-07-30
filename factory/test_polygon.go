package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// CreateTestTriangle creates a triangular polygon
func CreateTestTriangle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.PolygonCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent, components.GravityTag, components.ForceComponent)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 1.5) // Add some rotation

	// Create triangle vertices
	vertices := []Vec2.Vec2{
		{X: 0, Y: -50},  // top
		{X: -40, Y: 30}, // bottom-left
		{X: 40, Y: 30},  // bottom-right
	}

	poly := components.PolygonCollider.Get(entry)
	poly.Vertices = vertices

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 0.8, Y: 0.8}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 8
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a triangle (approximate)
	mc.Inertia = mc.Mass * 8000
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.7
	mat.StaticFriction = 0.3
	mat.DynamicFriction = 0.2

	return entry
}

// CreateTestPentagon creates a pentagonal polygon
func CreateTestPentagon(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.GravityTag, components.ForceComponent, components.MaterialComponent, components.Transform, components.PolygonCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 1.0) // Add some rotation

	// Create pentagon vertices
	vertices := components.CreateRegularPolygon(5, 60)

	poly := components.PolygonCollider.Get(entry)
	poly.Vertices = vertices

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 0.6, Y: 0.6}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 12
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a pentagon (approximate)
	mc.Inertia = mc.Mass * 12000
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.6
	mat.StaticFriction = 0.4
	mat.DynamicFriction = 0.3

	return entry
}

// CreateTestHexagon creates a hexagonal polygon
func CreateTestHexagon(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.GravityTag, components.ForceComponent, components.MaterialComponent, components.Transform, components.PolygonCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 0.8) // Add some rotation

	// Create hexagon vertices
	vertices := components.CreateRegularPolygon(6, 50)

	poly := components.PolygonCollider.Get(entry)
	poly.Vertices = vertices

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 0.7, Y: 0.7}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 10
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a hexagon (approximate)
	mc.Inertia = mc.Mass * 10000
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.75
	mat.StaticFriction = 0.35
	mat.DynamicFriction = 0.25

	return entry
}

// CreateTestRectangle creates a rectangular polygon (alternative to AABB)
func CreateTestRectangle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2, width, height float64) *donburi.Entry {
	entity := ecs.World.Create(components.GravityTag, components.ForceComponent, components.MaterialComponent, components.Transform, components.PolygonCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.TorqueComponent)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 1.2) // Add some rotation

	// Create rectangle vertices
	vertices := components.CreateRectangle(width, height)

	poly := components.PolygonCollider.Get(entry)
	poly.Vertices = vertices

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1.0, Y: 1.0}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 15
	mc.InverseMass = 1 / mc.Mass
	// Calculate inertia for a rectangle (approximate)
	mc.Inertia = mc.Mass * 15000
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.65
	mat.StaticFriction = 0.4
	mat.DynamicFriction = 0.3

	return entry
}
